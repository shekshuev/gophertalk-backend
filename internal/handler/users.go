package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	readDTOs, err := h.users.GetAllUsers()
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := json.Marshal(readDTOs)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}
