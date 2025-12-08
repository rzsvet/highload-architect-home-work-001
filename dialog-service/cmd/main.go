package main

import (
	"context"
	"dialog-service/internal/config"
	"dialog-service/internal/handler"
	"dialog-service/internal/repository"
	"dialog-service/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация MongoDB
	repo, err := repository.NewMongoDBRepository(cfg.MongoDBURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer repo.Close()

	// Инициализация сервисов
	dialogService := service.NewDialogService(repo)

	// Инициализация обработчиков
	dialogHandler := handler.NewDialogHandler(dialogService)

	// Настройка роутера
	router := gin.Default()

	// Middleware для логирования
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Routes
	router.GET("/health", dialogHandler.HealthCheck)
	router.POST("/dialog/:user_id/send", dialogHandler.SendMessage)
	router.GET("/dialog/:user_id/list", dialogHandler.GetDialog)

	// Запуск сервера
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
