package handlers

import (
	"encoding/json"
	"net/http"
)

func (s Server) CheckServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("ok")
}
