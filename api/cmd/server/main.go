package main

import (
	"api/internal/cache"
	"api/internal/config"
	"api/internal/database"
	"api/internal/handler"
	"api/internal/middleware"
	"api/internal/monitoring"
	"api/internal/repository"
	"api/internal/service"
	"fmt"
	"log"
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

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Continuing without cache...")
		// Можно создать заглушку или продолжить без кэша
		redisCache = nil
	} else {
		defer redisCache.Close()
	}

	// Initialize tables
	if err := database.InitTables(db); err != nil {
		log.Fatal("Failed to initialize tables:", err)
	}

	// Initialize services
	cacheService := service.NewCacheService(redisCache)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.WriteDB, db.ReadDB)
	friendRepo := repository.NewFriendRepository(db.WriteDB, db.ReadDB)
	postRepo := repository.NewPostRepository(db.WriteDB, db.ReadDB)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	friendService := service.NewFriendService(friendRepo, userRepo)
	// postService := service.NewPostService(postRepo)
	postService := service.NewPostService(postRepo, cacheService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	friendHandler := handler.NewFriendHandler(friendService)
	postHandler := handler.NewPostHandler(postService)
	searchHandler := handler.NewSearchHandler(userService)
	cacheHandler := handler.NewCacheHandler(cacheService, postService)

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
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}

	// Protected routes
	protected := router.Group(cfg.ServerPath)
	protected.Use(middleware.AuthMiddleware(userService))
	{

		// User routes
		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/users/:id", userHandler.GetUser)
		protected.GET("/profile", userHandler.GetProfile)
		// protected.PUT("/profile", userHandler.UpdateProfile)

		// Friend routes
		protected.POST("/friend/add", friendHandler.AddFriend)
		protected.POST("/friend/delete", friendHandler.DeleteFriend)
		protected.GET("/friends", friendHandler.GetFriends)
		protected.GET("/friend/status", friendHandler.GetFriendshipStatus)

		// Post routes
		protected.POST("/post/create", postHandler.CreatePost)
		protected.GET("/post/get/:id", postHandler.GetPost)
		protected.PUT("/post/update/:id", postHandler.UpdatePost)
		protected.DELETE("/post/delete/:id", postHandler.DeletePost)
		protected.GET("/posts", postHandler.GetUserPosts)
		protected.GET("/post/feed", postHandler.GetFeed)

		// Search routes
		protected.GET("/user/search", searchHandler.SearchUsers)
		protected.GET("/user/search/simple", searchHandler.SearchUsersSimple)

		// Сache routes
		protected.POST("/cache/invalidate", cacheHandler.InvalidateCache)
		protected.GET("/cache/stats", cacheHandler.GetCacheStats)

	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		writeDBErr := db.WriteDB.Ping()
		readDBErr := db.ReadDB.Ping()

		status := "ok"
		if writeDBErr != nil || readDBErr != nil {
			status = "degraded"
		}

		// Проверяем Redis если он подключен
		redisStatus := "not_configured"
		var redisStats map[string]interface{}

		if redisCache != nil {
			if err := redisCache.HealthCheck(); err != nil {
				redisStatus = "disconnected"
			} else {
				redisStatus = "connected"
				redisStats = redisCache.GetStats()
			}
		}

		response := gin.H{
			"status":   status,
			"write_db": map[string]interface{}{"connected": writeDBErr == nil},
			"read_db":  map[string]interface{}{"connected": readDBErr == nil},
			"redis":    redisStatus,
		}

		if redisStats != nil {
			response["redis_stats"] = redisStats
		}

		c.JSON(200, response)
	})

	// Start server

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if redisCache != nil {
		log.Printf("Redis cache: ENABLED")
	} else {
		log.Printf("Redis cache: DISABLED")
	}
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
