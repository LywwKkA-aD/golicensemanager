package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/LywwKkA-aD/golicensemanager/internal/app/handler"
	"github.com/LywwKkA-aD/golicensemanager/internal/config"
	"github.com/LywwKkA-aD/golicensemanager/internal/middleware"
	"github.com/LywwKkA-aD/golicensemanager/internal/repository/postgres"
	"github.com/LywwKkA-aD/golicensemanager/internal/service"
)

type App struct {
	config     *config.Config
	logger     *zap.SugaredLogger
	db         *gorm.DB
	router     *gin.Engine
	httpServer *http.Server
}

func NewApp(cfg *config.Config, logger *zap.SugaredLogger) (*App, error) {
	// Initialize database
	db, err := postgres.NewConnection(&postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := postgres.RunMigrations(db, "scripts/db/migrations"); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize repositories
	appRepo := postgres.NewApplicationRepository(db)
	licenseRepo := postgres.NewLicenseRepository(db)
	licenseTypeRepo := postgres.NewLicenseTypeRepository(db)
	clientRepo := postgres.NewClientRepository(db)

	// Initialize services
	appService := service.NewApplicationService(appRepo, logger)
	licenseService := service.NewLicenseService(licenseRepo, licenseTypeRepo, logger)
	clientService := service.NewClientService(clientRepo, licenseRepo, logger)

	// Initialize handlers
	appHandler := handler.NewApplicationHandler(appService, logger)
	licenseHandler := handler.NewLicenseHandler(licenseService, logger)
	clientHandler := handler.NewClientHandler(clientService, logger)

	// Initialize middlewares
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)
	corsMiddleware := middleware.NewCORSMiddleware(cfg.Server.AllowedOrigins)

	// Setup routes
	setupRoutes(router, *authMiddleware, *corsMiddleware, appHandler, licenseHandler, clientHandler)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:           cfg.Server.GetServerAddress(),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	return &App{
		config:     cfg,
		logger:     logger,
		db:         db,
		router:     router,
		httpServer: httpServer,
	}, nil
}

func (a *App) Start() error {
	// Start server in a goroutine
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	a.logger.Infof("Server is running on %s", a.config.Server.GetServerAddress())

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	// Close database connection
	if err := postgres.CloseConnection(a.db); err != nil {
		return fmt.Errorf("error closing database connection: %w", err)
	}

	return nil
}

func setupRoutes(
	r *gin.Engine,
	auth middleware.AuthMiddleware,
	cors middleware.CORSMiddleware,
	appHandler *handler.ApplicationHandler,
	licenseHandler *handler.LicenseHandler,
	clientHandler *handler.ClientHandler,
) {
	// Apply global middlewares
	r.Use(cors.Handler())
	r.Use(middleware.RequestLogger())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes
		v1.POST("/auth/token", appHandler.GenerateToken)

		// Protected routes
		authorized := v1.Group("")
		authorized.Use(auth.Handler())
		{
			// Application routes
			apps := authorized.Group("/applications")
			{
				apps.POST("", appHandler.Create)
				apps.GET("", appHandler.List)
				apps.GET("/:id", appHandler.Get)
				apps.PUT("/:id", appHandler.Update)
				apps.DELETE("/:id", appHandler.Delete)
			}

			// License routes
			licenses := authorized.Group("/licenses")
			{
				licenses.POST("", licenseHandler.Create)
				licenses.GET("", licenseHandler.List)
				licenses.GET("/:id", licenseHandler.Get)
				licenses.PUT("/:id", licenseHandler.Update)
				licenses.POST("/:id/revoke", licenseHandler.Revoke)
				licenses.POST("/:id/validate", licenseHandler.Validate)
			}

			// Client routes
			clients := authorized.Group("/clients")
			{
				clients.POST("", clientHandler.Create)
				clients.GET("", clientHandler.List)
				clients.GET("/:id", clientHandler.Get)
				clients.PUT("/:id", clientHandler.Update)
				clients.DELETE("/:id", clientHandler.Delete)
			}
		}
	}
}
