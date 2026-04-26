package main

import (
	"backend/internal/service/checklist_service"
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
	checklistRepo := repository.NewChecklistRepository(db)

	// JWT
	jwtManager := auth_service.NewJWTManager()

	// Services
	authService := auth_service.NewAuthService(userRepo, jwtManager)
	questionService := question_service.NewQuestionService(questionRepo)
	checklistService := checklist_service.NewChecklistService(checklistRepo, userRepo)
	reportService := report_service.NewReportService(
		db,
		reportRepo,
		questionRepo,
		checklistRepo,
		userRepo,
	)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	reportHandler := handler.NewReportHandler(reportService)
	questionHandler := handler.NewQuestionHandler(questionService)
	checklistHandler := handler.NewChecklistHandler(checklistService)

	// File storage
	fileStorage := storage.NewFileStorage(cfg.UploadDir, cfg.BaseURL)
	uploadHandler := handler.NewUploadHandler(fileStorage)

	// Router
	router := mux.NewRouter()

	// ===== API ROUTES (with /api prefix) =====
	api := router.PathPrefix("/api").Subrouter()

	// AUTH
	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")

	// CHECKLISTS
	api.HandleFunc("/checklists", checklistHandler.GetAll).Methods("GET")

	api.Handle("/checklists/available",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(checklistHandler.GetAvailableForCurrentUser),
		),
	).Methods("GET")

	api.Handle("/checklists/{id}/questions",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(questionHandler.GetByChecklist),
		),
	).Methods("GET")

	// UPLOAD
	api.Handle("/upload",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(uploadHandler.UploadImage),
		),
	).Methods("POST")

	// USER
	api.Handle("/me",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(authHandler.Me),
		),
	).Methods("GET")

	api.Handle("/admin",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("admin only"))
				}),
			),
		),
	).Methods("GET")

	// REPORTS
	api.Handle("/reports",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(reportHandler.GetReports),
			),
		),
	).Methods("GET")

	api.Handle("/reports",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(reportHandler.CreateReport),
		),
	).Methods("POST")

	api.Handle("/reports/export",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(reportHandler.ExportExcel),
			),
		),
	).Methods("GET")

	api.Handle("/reports/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(reportHandler.GetReportByID),
			),
		),
	).Methods("GET")

	// QUESTIONS
	api.Handle("/questions",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.GetAll),
			),
		),
	).Methods("GET")

	api.Handle("/questions",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.Create),
			),
		),
	).Methods("POST")

	api.Handle("/questions/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.Update),
			),
		),
	).Methods("PUT")

	api.Handle("/questions/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.Delete),
			),
		),
	).Methods("DELETE")

	// VARIETIES
	api.Handle("/varieties",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(varietyHandler.GetAll),
		),
	).Methods("GET")

	api.Handle("/varieties/{id}",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(varietyHandler.GetByID),
		),
	).Methods("GET")

	api.Handle("/varieties",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(varietyHandler.Create),
			),
		),
	).Methods("POST")

	api.Handle("/varieties/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(varietyHandler.Update),
			),
		),
	).Methods("PUT")

	api.Handle("/varieties/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(varietyHandler.Delete),
			),
		),
	).Methods("DELETE")

	// PHENOPHASES
	api.Handle("/phenophases",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(phenophaseHandler.GetAll),
		),
	).Methods("GET")

	api.Handle("/phenophases/{id}",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(phenophaseHandler.GetByID),
		),
	).Methods("GET")

	api.Handle("/phenophases",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(phenophaseHandler.Create),
			),
		),
	).Methods("POST")

	api.Handle("/phenophases/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(phenophaseHandler.Update),
			),
		),
	).Methods("PUT")

	api.Handle("/phenophases/{id}",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(phenophaseHandler.Delete),
			),
		),
	).Methods("DELETE")

	// STATIC FILES
	fs := http.FileServer(http.Dir(cfg.UploadDir))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))

	// CORS
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

	// SERVER
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handlerWithMiddleware,
	}

	go func() {
		log.Println("Server running on :" + cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

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
