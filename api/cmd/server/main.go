package main

import (
	"api/internal/config"
	"api/internal/database"
	"api/internal/handler"
	"api/internal/middleware"
	"api/internal/monitoring"
	"api/internal/repository"
	"api/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title User API Service
// @version 1.0
// @description API для регистрации, авторизации и управления пользователями
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		cfg := config.LoadConfig()
		db, err := database.NewDatabase(cfg)
		if err != nil {
			fmt.Println("FAILED: Cannot connect to database")
			os.Exit(1)
		}
		defer db.Close()

		if err := db.WriteDB.Ping(); err != nil {
			fmt.Println("FAILED: Write database not accessible")
			os.Exit(1)
		}

		if err := db.ReadDB.Ping(); err != nil {
			fmt.Println("FAILED: Read database not accessible")
			os.Exit(1)
		}

		fmt.Println("OK: All systems operational")
		os.Exit(0)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize tables
	if err := database.InitTables(db); err != nil {
		log.Fatal("Failed to initialize tables:", err)
	}

	// Initialize repository, service, and handler
	userRepo := repository.NewUserRepository(db.WriteDB, db.ReadDB)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userService)

	// Create Gin router
	router := gin.Default()

	router.Use(middleware.PrometheusMiddleware())

	if len(cfg.CORSAllowedOrigins) != 0 {
		router.Use(middleware.CORSWithConfig(cfg.CORSAllowedOrigins))
	}

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if cfg.ServerSwagger == "enabled" {
		// Swagger documentation
		swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL("/docs/openapi.yaml"),    // URL для загрузки манифеста
			ginSwagger.DefaultModelsExpandDepth(-1), // Скрываем модели по умолчанию
		)
		router.GET("/swagger/*any", swaggerHandler)

		// Также добавляем endpoint для raw JSON манифеста
		router.GET("/docs/openapi.yaml", func(c *gin.Context) {
			c.File("./docs/openapi.yaml") // Путь к сгенерированному файлу
		})
	}

	// Public routes
	public := router.Group(cfg.ServerPath)
	{
		public.POST("/user/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}

	// Protected routes
	protected := router.Group(cfg.ServerPath)
	protected.Use(middleware.AuthMiddleware(userService))
	{
		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/user/get/:id", userHandler.GetUser)
		protected.GET("/profile", userHandler.GetProfile)
		// protected.PUT("/profile", userHandler.UpdateProfile)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		// Проверяем соединение с базами данных
		writeDBErr := db.WriteDB.Ping()
		readDBErr := db.ReadDB.Ping()

		status := "ok"
		if writeDBErr != nil || readDBErr != nil {
			status = "degraded"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   status,
			"write_db": map[string]interface{}{"connected": writeDBErr == nil},
			"read_db":  map[string]interface{}{"connected": readDBErr == nil},
			"version":  "1.0.0",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("Metrics available at: http://localhost:%s/metrics", cfg.ServerPort)
	if cfg.ServerSwagger == "enabled" {
		log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", cfg.ServerPort)
	}
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// collectSystemMetrics собирает системные метрики
func collectSystemMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Количество goroutines
		monitoring.GoroutinesCount.Set(float64(runtime.NumGoroutine()))

		// Использование памяти
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		monitoring.MemoryUsage.WithLabelValues("alloc").Set(float64(m.Alloc))
		monitoring.MemoryUsage.WithLabelValues("sys").Set(float64(m.Sys))
		monitoring.MemoryUsage.WithLabelValues("heap_alloc").Set(float64(m.HeapAlloc))
		monitoring.MemoryUsage.WithLabelValues("heap_sys").Set(float64(m.HeapSys))
	}
}
