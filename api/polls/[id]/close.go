package api

import (
	"net/http"
	"net/url"
	"strings"
)

func Close(w http.ResponseWriter, r *http.Request) {
	if HandleCORS(w, r) {
		return
	}
	
	if r.Method != http.MethodPatch {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract poll ID from URL path
	urlPath, _ := url.PathUnescape(r.URL.Path)
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")
	
	if len(parts) < 4 || parts[0] != "api" || parts[1] != "polls" || parts[3] != "close" {
		WriteError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}
	
	pollID := parts[2]
	if pollID == "" {
		WriteError(w, http.StatusBadRequest, "Poll ID is required")
		return
	}
	
	store := GetStore()
	poll, err := store.ClosePoll(pollID)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, poll)
}