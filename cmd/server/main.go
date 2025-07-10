package main

import (
	"context"
	"gas_wells/internal/config"
	"gas_wells/internal/handler"
	"gas_wells/internal/pkg/database"
	"gas_wells/internal/pkg/logger"
	"gas_wells/internal/repository"
	"gas_wells/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Инициализация логгера
	log := logger.New(cfg.App.Env)

	// Инициализация базы данных
	db, err := database.NewPostgres(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	})

	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Инициализация репозитория
	wellRepo := repository.NewWellRepo(db.Pool, log)

	// Инициализация сервиса
	wellService := service.NewWellService(wellRepo, log)

	// Инициализация обработчика
	wellHandler := handler.NewWellHandler(wellService, log)

	// Настройка сервера
	router := setupRouter(wellHandler, log)

	// HTTP сервер
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Graceful shutdown
	go func() {
		log.Info("starting server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown error", "error", err)
	}

	log.Info("server stopped gracefully")
}

func setupRouter(wellHandler *handler.WellHandler, log logger.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(loggerMiddleware(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Маршруты
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Gas Wells API"))
	})

	r.Route("/wells", func(r chi.Router) {
		r.Get("/", wellHandler.ListWells)
		r.Post("/", wellHandler.CreateWell)
		r.Get("/{id}", wellHandler.GetWell)
		r.Put("/{id}", wellHandler.UpdateWell)
		r.Delete("/{id}", wellHandler.DeleteWell)
	})

	return r
}

func loggerMiddleware(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				log.Info("request completed",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"duration", time.Since(start))
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
