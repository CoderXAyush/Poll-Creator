package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Polls(w http.ResponseWriter, r *http.Request) {
	if HandleCORS(w, r) {
		return
	}
	
	store := GetStore()
	
	switch r.Method {
	case http.MethodGet:
		polls := store.GetAllPolls()
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
		WriteJSON(w, http.StatusOK, items)
		
	case http.MethodPost:
		var req CreatePollRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if strings.TrimSpace(req.Question) == "" || len(req.Options) < 2 {
			WriteError(w, http.StatusBadRequest, "Question and at least 2 options are required.")
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
			WriteError(w, http.StatusBadRequest, "At least 2 non-empty options are required.")
			return
		}

		poll := store.CreatePoll(req.Question, validOpts)
		WriteJSON(w, http.StatusCreated, poll)
		
	default:
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}