package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service/auth_service"
	"backend/internal/service/question_service"
	"backend/internal/service/report_service"
	"backend/internal/storage"
)

func main() {
	cfg := config.LoadConfig()

	// DB
	db := repository.NewDB(cfg)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	reportRepo := repository.NewReportRepository(db)
	questionRepo := repository.NewQuestionRepository(db)

	// JWT
	jwtManager := auth_service.NewJWTManager()

	// Services
	authService := auth_service.NewAuthService(userRepo, jwtManager)
	reportService := report_service.NewReportService(db, reportRepo)
	questionService := question_service.NewQuestionService(questionRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	reportHandler := handler.NewReportHandler(reportService)
	questionHandler := handler.NewQuestionHandler(questionService)

	// File storage
	fileStorage := storage.NewFileStorage(cfg.UploadDir, cfg.BaseURL)
	uploadHandler := handler.NewUploadHandler(fileStorage)

	// Router
	router := mux.NewRouter()

	// ===== AUTH =====
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	// ===== UPLOAD =====
	router.Handle("/upload",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(uploadHandler.UploadImage),
		),
	).Methods("POST")

	// ===== TEST =====
	router.Handle("/me",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value(middleware.UserIDKey)
				role := r.Context().Value(middleware.RoleKey)

				w.Write([]byte("Hello user " + userID.(string) + " role: " + role.(string)))
			}),
		),
	).Methods("GET")

	router.Handle("/admin",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("admin only"))
				}),
			),
		),
	).Methods("GET")

	// ===== REPORTS =====
	router.Handle("/reports",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(reportHandler.GetReports),
			),
		),
	).Methods("GET")

	router.Handle("/reports",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(reportHandler.CreateReport),
		),
	).Methods("POST")

	router.Handle("/reports/export",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(reportHandler.ExportExcel),
			),
		),
	).Methods("GET")

	router.Handle("/reports/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(reportHandler.GetReportByID),
			),
		),
	).Methods("GET")

	// ===== QUESTIONS =====
	router.Handle("/questions",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(questionHandler.GetAll),
		),
	).Methods("GET")

	router.Handle("/questions",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.Create),
			),
		),
	).Methods("POST")

	router.Handle("/questions/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.Update),
			),
		),
	).Methods("PUT")

	router.Handle("/questions/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.Delete),
			),
		),
	).Methods("DELETE")

	// ===== STATIC FILES (uploads) =====
	fs := http.FileServer(http.Dir(cfg.UploadDir))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))

	// ===== MIDDLEWARE =====
	allowedOrigins := make(map[string]bool)
	for _, origin := range cfg.CorsAllowedOrigins {
		if origin != "" {
			allowedOrigins[origin] = true
		}
	}

	corsMiddleware := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: cfg.CorsAllowCredentials,
	})
	handlerWithMiddleware := corsMiddleware(router)

	// ===== SERVER =====
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handlerWithMiddleware,
	}

	// ===== START SERVER =====
	go func() {
		log.Println("Server running on :" + cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// ===== GRACEFUL SHUTDOWN =====
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
