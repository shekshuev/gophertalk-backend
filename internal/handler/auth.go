package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/shekshuev/gophertalk-backend/internal/models"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO models.LoginUserDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &loginDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokensDTO, err := h.auth.Login(loginDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(tokensDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
