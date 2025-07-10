package main

import (
	"context"
	"gas_wells/internal/handler"
	"gas_wells/internal/repository"
	"gas_wells/internal/server"
	"gas_wells/internal/service"
	"gas_wells/pkg/database"
	"gas_wells/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Инициализация логгера
	log := logger.New(os.Getenv("APP_ENV"))

	// Инициализация БД
	db, err := database.NewPostgres(database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "gas_wells"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	})
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Применение миграций
	migrator := database.NewMigrator(db.Pool, log)
	if err := migrator.Run(context.Background()); err != nil {
		log.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}

	// Инициализация приложения
	wellRepo := repository.NewWellRepo(db.Pool, log)
	wellService := service.NewWellService(wellRepo, log)
	wellHandler := handler.NewWellHandler(wellService, log)

	// Настройка сервера
	srv := server.New(log)
	srv.SetupRoutes(wellHandler)
	srv.ServeStatic("web/static")

	// Запуск сервера
	httpSrv := &http.Server{
		Addr:         ":" + getEnv("PORT", "8080"),
		Handler:      srv.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Starting server", "port", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", "error", err)
	}

	log.Info("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
