package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	readDTOs, err := h.users.GetAllUsers(limit, offset)
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

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}
	readDTO, err := h.users.GetUserByID(id)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	resp, err := json.Marshal(readDTO)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}
