package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ─── Models ────────────────────────────────────────────────────────────────────

type Option struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

type Poll struct {
	ID         string   `json:"id"`
	Question   string   `json:"question"`
	Options    []Option `json:"options"`
	TotalVotes int      `json:"totalVotes"`
	Closed     bool     `json:"closed"`
	CreatedAt  string   `json:"createdAt"`
}

type PollResponse struct {
	Poll
	HasVoted      bool  `json:"hasVoted"`
	VotedOptionID *int  `json:"votedOptionId"`
}

type PollListItem struct {
	ID          string `json:"id"`
	Question    string `json:"question"`
	TotalVotes  int    `json:"totalVotes"`
	Closed      bool   `json:"closed"`
	CreatedAt   string `json:"createdAt"`
	OptionCount int    `json:"optionCount"`
}

type ResultOption struct {
	ID         int     `json:"id"`
	Text       string  `json:"text"`
	Votes      int     `json:"votes"`
	Percentage float64 `json:"percentage"`
}

type PollResults struct {
	ID         string         `json:"id"`
	Question   string         `json:"question"`
	Options    []ResultOption `json:"options"`
	TotalVotes int            `json:"totalVotes"`
	Closed     bool           `json:"closed"`
	CreatedAt  string         `json:"createdAt"`
}

type CreatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type VoteRequest struct {
	OptionID *int `json:"optionId"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// ─── Store ─────────────────────────────────────────────────────────────────────

type Store struct {
	mu    sync.RWMutex
	polls map[string]*Poll
	votes map[string]map[string]int // pollID -> voterKey -> optionID
}

func NewStore() *Store {
	return &Store{
		polls: make(map[string]*Poll),
		votes: make(map[string]map[string]int),
	}
}

func (s *Store) CreatePoll(question string, options []string) *Poll {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	opts := make([]Option, len(options))
	for i, text := range options {
		opts[i] = Option{ID: i, Text: strings.TrimSpace(text), Votes: 0}
	}

	poll := &Poll{
		ID:         id,
		Question:   strings.TrimSpace(question),
		Options:    opts,
		TotalVotes: 0,
		Closed:     false,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	s.polls[id] = poll
	s.votes[id] = make(map[string]int)
	return poll
}

func (s *Store) GetPoll(id string) *Poll {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.polls[id]
}

func (s *Store) GetAllPolls() []*Poll {
	s.mu.RLock()
	defer s.mu.RUnlock()

	polls := make([]*Poll, 0, len(s.polls))
	for _, p := range s.polls {
		polls = append(polls, p)
	}
	return polls
}

func (s *Store) HasVoted(pollID, voterKey string) (bool, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if voters, ok := s.votes[pollID]; ok {
		if optID, voted := voters[voterKey]; voted {
			return true, optID
		}
	}
	return false, -1
}

func (s *Store) Vote(pollID, voterKey string, optionID int) (*Poll, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	poll, ok := s.polls[pollID]
	if !ok {
		return nil, fmt.Errorf("Poll not found")
	}
	if poll.Closed {
		return nil, fmt.Errorf("This poll is closed")
	}

	// Validate option
	if optionID < 0 || optionID >= len(poll.Options) {
		return nil, fmt.Errorf("Invalid option")
	}

	// Check duplicate
	if _, voted := s.votes[pollID][voterKey]; voted {
		return nil, fmt.Errorf("You have already voted on this poll")
	}

	poll.Options[optionID].Votes++
	poll.TotalVotes++
	s.votes[pollID][voterKey] = optionID

	return poll, nil
}

func (s *Store) ClosePoll(id string) (*Poll, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	poll, ok := s.polls[id]
	if !ok {
		return nil, fmt.Errorf("Poll not found")
	}
	poll.Closed = true
	return poll, nil
}

// ─── Seed ──────────────────────────────────────────────────────────────────────

func (s *Store) SeedDefault() {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	poll := &Poll{
		ID:       id,
		Question: "What's your favorite programming language?",
		Options: []Option{
			{ID: 0, Text: "JavaScript", Votes: 12},
			{ID: 1, Text: "Python", Votes: 18},
			{ID: 2, Text: "Go", Votes: 25},
			{ID: 3, Text: "Rust", Votes: 8},
			{ID: 4, Text: "TypeScript", Votes: 15},
			{ID: 5, Text: "Java", Votes: 10},
		},
		TotalVotes: 88,
		Closed:     false,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	s.polls[id] = poll
	s.votes[id] = make(map[string]int)
	log.Printf("Seeded default poll: %s", id)
}

// ─── Helpers ───────────────────────────────────────────────────────────────────

func getVoterKey(r *http.Request) string {
	if sid := r.Header.Get("X-Session-Id"); sid != "" {
		return sid
	}
	// Use X-Forwarded-For or RemoteAddr
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

// ─── CORS Middleware ───────────────────────────────────────────────────────────

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-Id")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ─── Handlers ──────────────────────────────────────────────────────────────────

type Server struct {
	store *Store
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleCreatePoll(w http.ResponseWriter, r *http.Request) {
	var req CreatePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Question) == "" || len(req.Options) < 2 {
		writeError(w, http.StatusBadRequest, "Question and at least 2 options are required.")
		return
	}

	// Filter empty options
	var validOpts []string
	for _, o := range req.Options {
		if strings.TrimSpace(o) != "" {
			validOpts = append(validOpts, o)
		}
	}
	if len(validOpts) < 2 {
		writeError(w, http.StatusBadRequest, "At least 2 non-empty options are required.")
		return
	}

	poll := s.store.CreatePoll(req.Question, validOpts)
	writeJSON(w, http.StatusCreated, poll)
}

func (s *Server) handleGetPolls(w http.ResponseWriter, r *http.Request) {
	polls := s.store.GetAllPolls()
	items := make([]PollListItem, len(polls))
	for i, p := range polls {
		items[i] = PollListItem{
			ID:          p.ID,
			Question:    p.Question,
			TotalVotes:  p.TotalVotes,
			Closed:      p.Closed,
			CreatedAt:   p.CreatedAt,
			OptionCount: len(p.Options),
		}
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleGetPoll(w http.ResponseWriter, r *http.Request, pollID string) {
	poll := s.store.GetPoll(pollID)
	if poll == nil {
		writeError(w, http.StatusNotFound, "Poll not found.")
		return
	}

	voterKey := getVoterKey(r)
	hasVoted, votedOptID := s.store.HasVoted(pollID, voterKey)

	resp := PollResponse{Poll: *poll, HasVoted: hasVoted}
	if hasVoted {
		resp.VotedOptionID = &votedOptID
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleVote(w http.ResponseWriter, r *http.Request, pollID string) {
	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.OptionID == nil {
		writeError(w, http.StatusBadRequest, "optionId is required.")
		return
	}

	voterKey := getVoterKey(r)
	poll, err := s.store.Vote(pollID, voterKey, *req.OptionID)
	if err != nil {
		msg := err.Error()
		status := http.StatusBadRequest
		if msg == "Poll not found" {
			status = http.StatusNotFound
		} else if msg == "You have already voted on this poll" {
			status = http.StatusConflict
		}
		writeError(w, status, msg)
		return
	}

	resp := PollResponse{Poll: *poll, HasVoted: true, VotedOptionID: req.OptionID}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleGetResults(w http.ResponseWriter, r *http.Request, pollID string) {
	poll := s.store.GetPoll(pollID)
	if poll == nil {
		writeError(w, http.StatusNotFound, "Poll not found.")
		return
	}

	results := PollResults{
		ID:         poll.ID,
		Question:   poll.Question,
		TotalVotes: poll.TotalVotes,
		Closed:     poll.Closed,
		CreatedAt:  poll.CreatedAt,
	}

	results.Options = make([]ResultOption, len(poll.Options))
	for i, o := range poll.Options {
		pct := 0.0
		if poll.TotalVotes > 0 {
			pct = math.Round(float64(o.Votes)/float64(poll.TotalVotes)*1000) / 10
		}
		results.Options[i] = ResultOption{
			ID:         o.ID,
			Text:       o.Text,
			Votes:      o.Votes,
			Percentage: pct,
		}
	}

	writeJSON(w, http.StatusOK, results)
}

func (s *Server) handleClosePoll(w http.ResponseWriter, r *http.Request, pollID string) {
	poll, err := s.store.ClosePoll(pollID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, poll)
}

// ─── Router ────────────────────────────────────────────────────────────────────

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	// API routes
	if strings.HasPrefix(path, "/api") {
		s.handleAPIRoutes(w, r, path)
		return
	}

	// Static file serving for frontend
	s.handleStaticFiles(w, r)
}

func (s *Server) handleAPIRoutes(w http.ResponseWriter, r *http.Request, path string) {
	// GET /api/health
	if path == "/api/health" && r.Method == http.MethodGet {
		s.handleHealth(w, r)
		return
	}

	// POST /api/polls  |  GET /api/polls
	if path == "/api/polls" {
		switch r.Method {
		case http.MethodPost:
			s.handleCreatePoll(w, r)
		case http.MethodGet:
			s.handleGetPolls(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// Routes with poll ID: /api/polls/{id}[/vote|/results|/close]
	if strings.HasPrefix(path, "/api/polls/") {
		rest := strings.TrimPrefix(path, "/api/polls/")
		parts := strings.SplitN(rest, "/", 2)
		pollID := parts[0]

		if len(parts) == 1 {
			// GET /api/polls/:id
			if r.Method == http.MethodGet {
				s.handleGetPoll(w, r, pollID)
				return
			}
		} else {
			sub := parts[1]
			switch {
			case sub == "vote" && r.Method == http.MethodPost:
				s.handleVote(w, r, pollID)
				return
			case sub == "results" && r.Method == http.MethodGet:
				s.handleGetResults(w, r, pollID)
				return
			case sub == "close" && r.Method == http.MethodPatch:
				s.handleClosePoll(w, r, pollID)
				return
			}
		}
	}

	writeError(w, http.StatusNotFound, "API endpoint not found")
}

func (s *Server) handleStaticFiles(w http.ResponseWriter, r *http.Request) {
	// Check if static directory exists
	staticDir := "./static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		// If static directory doesn't exist, show a simple message
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Poll Creator API</title>
			</head>
			<body>
				<h1>🗳️ Poll Creator API</h1>
				<p>Backend API is running successfully!</p>
				<p>API endpoints are available at <code>/api/*</code></p>
				<ul>
					<li><a href="/api/health">Health Check</a></li>
					<li><a href="/api/polls">List Polls</a></li>
				</ul>
			</body>
			</html>
		`)
		return
	}

	// Serve static files
	requestPath := r.URL.Path
	if requestPath == "/" {
		requestPath = "/index.html"
	}

	filePath := filepath.Join(staticDir, requestPath)

	// Security check: ensure the file is within staticDir
	if !strings.HasPrefix(filePath, staticDir+string(os.PathSeparator)) && filePath != staticDir+"/index.html" {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// For SPA routing, serve index.html for non-API routes
		indexPath := filepath.Join(staticDir, "index.html")
		if _, indexErr := os.Stat(indexPath); indexErr == nil {
			http.ServeFile(w, r, indexPath)
			return
		}
		writeError(w, http.StatusNotFound, "Not found")
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// ─── Main ──────────────────────────────────────────────────────────────────────

func main() {
	store := NewStore()
	store.SeedDefault()

	server := &Server{store: store}

	// Get port from environment variable, default to 5000
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 Poll Creator API (Go) running on http://0.0.0.0:%s", port)
	log.Printf("📁 Static files will be served from ./static (if available)")
	log.Printf("🔧 Environment: %s", getEnv())
	
	if err := http.ListenAndServe(":"+port, corsMiddleware(server)); err != nil {
		log.Fatal(err)
	}
}

func getEnv() string {
	if env := os.Getenv("NODE_ENV"); env != "" {
		return env
	}
	if os.Getenv("PORT") != "" {
		return "production"
	}
	return "development"
}
