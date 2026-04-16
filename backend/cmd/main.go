package main

import (
	"backend/internal/middleware"
	"backend/internal/service/question_service"
	"backend/internal/service/report_service"
	"backend/internal/storage"
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

	questionRepo := repository.NewQuestionRepository(db)
	questionService := question_service.NewQuestionService(questionRepo)
	questionHandler := handler.NewQuestionHandler(questionService)

	mux := http.NewServeMux()

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	protected := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey)
		role := r.Context().Value(middleware.RoleKey)

		w.Write([]byte("Hello user " + userID.(string) + " role: " + role.(string)))
	})

	fs := http.FileServer(http.Dir(cfg.UploadDir))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	fileStorage := storage.NewFileStorage(cfg.UploadDir, cfg.BaseURL)

	uploadHandler := handler.NewUploadHandler(fileStorage)

	mux.Handle("/upload",
		middleware.AuthMiddleware(jwtManager)(
			http.HandlerFunc(uploadHandler.UploadImage),
		),
	)

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

	mux.Handle("/get-questions",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.GetAll),
			),
		),
	)

	mux.Handle("/create-questions",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.Create),
			),
		),
	)

	mux.Handle("/update-questions/",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.Update),
			),
		),
	)

	mux.Handle("/delete-questions/",
		middleware.AuthMiddleware(jwtManager)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(questionHandler.Delete),
			),
		),
	)

	log.Println("Server running on :" + cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, mux)
}
