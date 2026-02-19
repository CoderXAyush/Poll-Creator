package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
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
	HasVoted      bool `json:"hasVoted"`
	VotedOptionID *int `json:"votedOptionId"`
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

// ─── Store (in-memory, shared across warm invocations) ─────────────────────────

type Store struct {
	polls map[string]*Poll
	votes map[string]map[string]int
	mu    sync.RWMutex
}

var store *Store
var once sync.Once

func getStore() *Store {
	once.Do(func() {
		store = &Store{
			polls: make(map[string]*Poll),
			votes: make(map[string]map[string]int),
		}
		seedDefault(store)
	})
	return store
}

func seedDefault(s *Store) {
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
}

// ─── Store Methods ─────────────────────────────────────────────────────────────

func (s *Store) CreatePoll(question string, options []string) *Poll {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	opts := make([]Option, len(options))
	for i, opt := range options {
		opts[i] = Option{ID: i, Text: strings.TrimSpace(opt), Votes: 0}
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
	if optionID < 0 || optionID >= len(poll.Options) {
		return nil, fmt.Errorf("Invalid option")
	}
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

// ─── Helpers ───────────────────────────────────────────────────────────────────

func getVoterKey(r *http.Request) string {
	if sid := r.Header.Get("X-Session-Id"); sid != "" {
		return sid
	}
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

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-Id")
}

// ─── Main Handler (Vercel entrypoint) ──────────────────────────────────────────

func Handler(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	s := getStore()
	path := strings.TrimSuffix(r.URL.Path, "/")

	// GET /api/health
	if path == "/api/health" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	// POST /api/polls  |  GET /api/polls
	if path == "/api/polls" {
		switch r.Method {
		case http.MethodGet:
			handleGetPolls(w, s)
		case http.MethodPost:
			handleCreatePoll(w, r, s)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// Routes: /api/polls/{id}[/vote|/results|/close]
	if strings.HasPrefix(path, "/api/polls/") {
		rest := strings.TrimPrefix(path, "/api/polls/")
		parts := strings.SplitN(rest, "/", 2)
		pollID := parts[0]

		if len(parts) == 1 {
			if r.Method == http.MethodGet {
				handleGetPoll(w, r, s, pollID)
				return
			}
		} else {
			sub := parts[1]
			switch {
			case sub == "vote" && r.Method == http.MethodPost:
				handleVote(w, r, s, pollID)
				return
			case sub == "results" && r.Method == http.MethodGet:
				handleGetResults(w, s, pollID)
				return
			case sub == "close" && r.Method == http.MethodPatch:
				handleClosePoll(w, s, pollID)
				return
			}
		}
	}

	writeError(w, http.StatusNotFound, "Not found")
}

// ─── Route Handlers ────────────────────────────────────────────────────────────

func handleGetPolls(w http.ResponseWriter, s *Store) {
	polls := s.GetAllPolls()
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

func handleCreatePoll(w http.ResponseWriter, r *http.Request, s *Store) {
	var req CreatePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Question) == "" || len(req.Options) < 2 {
		writeError(w, http.StatusBadRequest, "Question and at least 2 options are required.")
		return
	}

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

	poll := s.CreatePoll(req.Question, validOpts)
	writeJSON(w, http.StatusCreated, poll)
}

func handleGetPoll(w http.ResponseWriter, r *http.Request, s *Store, pollID string) {
	poll := s.GetPoll(pollID)
	if poll == nil {
		writeError(w, http.StatusNotFound, "Poll not found.")
		return
	}

	voterKey := getVoterKey(r)
	hasVoted, votedOptID := s.HasVoted(pollID, voterKey)

	resp := PollResponse{Poll: *poll, HasVoted: hasVoted}
	if hasVoted {
		resp.VotedOptionID = &votedOptID
	}
	writeJSON(w, http.StatusOK, resp)
}

func handleVote(w http.ResponseWriter, r *http.Request, s *Store, pollID string) {
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
	poll, err := s.Vote(pollID, voterKey, *req.OptionID)
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

func handleGetResults(w http.ResponseWriter, s *Store, pollID string) {
	poll := s.GetPoll(pollID)
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

func handleClosePoll(w http.ResponseWriter, s *Store, pollID string) {
	poll, err := s.ClosePoll(pollID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, poll)
}
