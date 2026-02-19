package api

import (
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	if HandleCORS(w, r) {
		return
	}
	
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}