package api

import (
	"math"
	"net/http"
	"net/url"
	"strings"
)

func Results(w http.ResponseWriter, r *http.Request) {
	if HandleCORS(w, r) {
		return
	}
	
	// Extract poll ID from URL path
	urlPath, _ := url.PathUnescape(r.URL.Path)
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")
	
	if len(parts) < 4 || parts[0] != "api" || parts[1] != "polls" || parts[3] != "results" {
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

	WriteJSON(w, http.StatusOK, results)
}