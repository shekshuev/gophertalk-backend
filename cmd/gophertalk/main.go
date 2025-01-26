package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/handler"
	"github.com/shekshuev/gophertalk-backend/internal/repository"
	"github.com/shekshuev/gophertalk-backend/internal/service"
)

func main() {
	cfg := config.GetConfig()
	userRepo := repository.NewUserRepositoryImpl(&cfg)
	postRepo := repository.NewPostRepositoryImpl(&cfg)
	userService := service.NewUserServiceImpl(userRepo, &cfg)
	authService := service.NewAuthServiceImpl(userRepo, &cfg)
	postService := service.NewPostServiceImpl(postRepo, &cfg)
	userHandler := handler.NewHandler(userService, authService, postService, &cfg)
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: userHandler.Router,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error starting server")
		}
	}()
	log.Print("Server listening on ", cfg.ServerAddress)
	<-done
	log.Print("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown")
	}
	log.Print("Server shutdown gracefully")
}
