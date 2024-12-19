package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	readDTOs, err := h.users.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	resp, err := json.Marshal(readDTOs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
