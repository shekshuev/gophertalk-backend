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
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err = json.Unmarshal(body, &loginDTO); err != nil {
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.validate.Struct(loginDTO)
	if err != nil {
		h.JSONError(w, http.StatusUnprocessableEntity, ErrValidationError.Error())
		return
	}
	tokensDTO, err := h.auth.Login(loginDTO)
	if err != nil {
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	resp, err := json.Marshal(tokensDTO)
	if err != nil {
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusUnauthorized, err.Error())
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO models.RegisterUserDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err = json.Unmarshal(body, &registerDTO); err != nil {
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.validate.Struct(registerDTO)
	if err != nil {
		h.JSONError(w, http.StatusUnprocessableEntity, ErrValidationError.Error())
		return
	}
	tokensDTO, err := h.auth.Register(registerDTO)
	if err != nil {
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	resp, err := json.Marshal(tokensDTO)
	if err != nil {
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusUnauthorized, err.Error())
	}
}
