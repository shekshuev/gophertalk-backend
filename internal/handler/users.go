package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/shekshuev/gophertalk-backend/internal/service"
)

type UserHandler struct {
	service service.UserService
	Router  *chi.Mux
}

func NewUserHandler(service service.UserService) *UserHandler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)
	h := &UserHandler{service: service, Router: router}
	h.Router.Get("/v1.0/users", h.GetAllUsers)
	return h
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	readDTOs, err := h.service.GetAllUsers()
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
