package api

import (
	"net/http"
	"net/url"
	"strings"
)

func PollDetails(w http.ResponseWriter, r *http.Request) {
	if HandleCORS(w, r) {
		return
	}
	
	// Extract poll ID from URL path
	urlPath, _ := url.PathUnescape(r.URL.Path)
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")
	
	if len(parts) < 3 || parts[0] != "api" || parts[1] != "polls" {
		WriteError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}
	
	pollID := parts[2]
	if pollID == "" {
		WriteError(w, http.StatusBadRequest, "Poll ID is required")
		return
	}
	
	store := GetStore()
	poll := store.GetPoll(pollID)
	if poll == nil {
		WriteError(w, http.StatusNotFound, "Poll not found.")
		return
	}

	voterKey := GetVoterKey(r)
	hasVoted, votedOptID := store.HasVoted(pollID, voterKey)

	resp := PollResponse{Poll: *poll, HasVoted: hasVoted}
	if hasVoted {
		resp.VotedOptionID = &votedOptID
	}

	WriteJSON(w, http.StatusOK, resp)
}