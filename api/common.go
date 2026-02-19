package api

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
	polls map[string]*Poll
	votes map[string]map[string]int // pollID -> voterKey -> optionID
	mu    sync.RWMutex
}

var globalStore *Store
var once sync.Once

func GetStore() *Store {
	once.Do(func() {
		globalStore = &Store{
			polls: make(map[string]*Poll),
			votes: make(map[string]map[string]int),
		}
		globalStore.SeedDefault()
	})
	return globalStore
}

func (s *Store) CreatePoll(question string, options []string) *Poll {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	opts := make([]Option, len(options))
	for i, opt := range options {
		opts[i] = Option{
			ID:    i,
			Text:  strings.TrimSpace(opt),
			Votes: 0,
		}
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
}

// ─── Helpers ───────────────────────────────────────────────────────────────────

func GetVoterKey(r *http.Request) string {
	if sid := r.Header.Get("X-Session-Id"); sid != "" {
		return sid
	}
	// Use X-Forwarded-For or RemoteAddr
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-Id")
	
	if status != 0 {
		w.WriteHeader(status)
	}
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
}

func HandleCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-Id")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}