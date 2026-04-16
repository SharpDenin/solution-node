package main

import (
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/service/auth"
)

func main() {
	cfg := config.LoadConfig()

	// DB
	db := repository.NewDB(cfg)

	// repo
	userRepo := repository.NewUserRepository(db)

	// jwt
	jwtManager := auth.NewJWTManager()

	// service
	authService := auth.NewAuthService(userRepo, jwtManager)

	// handler
	authHandler := handlers.NewAuthHandler(authService)

	// routes
	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)

	log.Println("Server running on :" + cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, nil)
}
