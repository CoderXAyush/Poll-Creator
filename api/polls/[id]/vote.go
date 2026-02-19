package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func Vote(w http.ResponseWriter, r *http.Request) {
	if HandleCORS(w, r) {
		return
	}
	
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract poll ID from URL path
	urlPath, _ := url.PathUnescape(r.URL.Path)
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")
	
	if len(parts) < 4 || parts[0] != "api" || parts[1] != "polls" || parts[3] != "vote" {
		WriteError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}
	
	pollID := parts[2]
	if pollID == "" {
		WriteError(w, http.StatusBadRequest, "Poll ID is required")
		return
	}
	
	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.OptionID == nil {
		WriteError(w, http.StatusBadRequest, "optionId is required.")
		return
	}

	store := GetStore()
	voterKey := GetVoterKey(r)
	poll, err := store.Vote(pollID, voterKey, *req.OptionID)
	if err != nil {
		msg := err.Error()
		status := http.StatusBadRequest
		if msg == "Poll not found" {
			status = http.StatusNotFound
		} else if msg == "You have already voted on this poll" {
			status = http.StatusConflict
		}
		WriteError(w, status, msg)
		return
	}

	resp := PollResponse{Poll: *poll, HasVoted: true, VotedOptionID: req.OptionID}
	WriteJSON(w, http.StatusOK, resp)
}