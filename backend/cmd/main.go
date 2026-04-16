package main

import (
	"backend/internal/middleware"
	"backend/internal/service/report_service"
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service/auth_service"
)

func main() {
	cfg := config.LoadConfig()

	db := repository.NewDB(cfg)
	userRepo := repository.NewUserRepository(db)

	jwtManager := auth_service.NewJWTManager()
	authService := auth_service.NewAuthService(userRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	reportRepo := repository.NewReportRepository(db)
	reportService := report_service.NewReportService(db, reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

	mux := http.NewServeMux()

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	protected := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey)
		role := r.Context().Value(middleware.RoleKey)

		w.Write([]byte("Hello user " + userID.(string) + " role: " + role.(string)))
	})

	mux.Handle("/me",
		middleware.AuthMiddleware(jwtManager)(
			protected,
		),
	)

	mux.Handle("/admin",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("admin only"))
				}),
			),
		),
	)

	mux.Handle("/create-report",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(reportHandler.CreateReport),
		),
	)

	mux.Handle("/get-report-list",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(reportHandler.GetReports),
			),
		),
	)

	mux.Handle("/reports/",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(reportHandler.GetReportByID),
			),
		),
	)

	log.Println("Server running on :" + cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, mux)
}
