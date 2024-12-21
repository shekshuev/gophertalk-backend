package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/middleware"
	"github.com/shekshuev/gophertalk-backend/internal/service"
	"github.com/shekshuev/gophertalk-backend/internal/utils"
)

type Handler struct {
	users    service.UserService
	auth     service.AuthService
	posts    service.PostService
	Router   *chi.Mux
	validate *validator.Validate
	cfg      *config.Config
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var ErrValidationError = errors.New("validation error")
var ErrInvalidID = errors.New("invalid ID")
var ErrInvalidToken = errors.New("invalid token")

func NewHandler(
	users service.UserService,
	auth service.AuthService,
	posts service.PostService,
	cfg *config.Config,
) *Handler {
	router := chi.NewRouter()
	validate := utils.NewValidator()
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))
	router.Use(chiMiddleware.Recoverer)
	router.Use(cors.AllowAll().Handler)
	h := &Handler{users: users, auth: auth, posts: posts, Router: router, validate: validate, cfg: cfg}

	h.Router.Route("/v1.0/users", func(r chi.Router) {
		r.Use(middleware.RequestAuth(cfg.AccessTokenSecret))
		r.Get("/", h.GetAllUsers)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetUserByID)
			r.With(middleware.RequestAuthSameID(cfg.AccessTokenSecret)).Put("/", h.UpdateUser)
			r.With(middleware.RequestAuthSameID(cfg.AccessTokenSecret)).Delete("/", h.DeleteUserByID)
		})
	})

	h.Router.Route("/v1.0/posts", func(r chi.Router) {
		r.Use(middleware.RequestAuth(cfg.AccessTokenSecret))
		r.Get("/", h.GetAllPosts)
		r.Post("/", h.CreatePost)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetPostByID)
			r.Delete("/", h.DeletePostByID)
			r.Post("/view", h.ViewPost)
			r.Post("/like", h.LikePost)
			r.Delete("/like", h.DislikePost)
		})
	})

	h.Router.Route("/v1.0/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
	})

	return h
}

func (h *Handler) JSONError(w http.ResponseWriter, statusCode int, errMessage string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errMessage})
}
