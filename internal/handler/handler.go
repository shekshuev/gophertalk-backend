package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/shekshuev/gophertalk-backend/internal/service"
)

type Handler struct {
	users  service.UserService
	auth   service.AuthService
	Router *chi.Mux
}

func NewHandler(users service.UserService, auth service.AuthService) *Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)
	h := &Handler{users: users, auth: auth, Router: router}

	h.Router.Route("/v1.0/users", func(r chi.Router) {
		r.Get("/", h.GetAllUsers)
	})

	h.Router.Route("/v1.0/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
	})

	return h
}
