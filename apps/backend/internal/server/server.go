package server

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/visstt/multisub/backend/internal/auth"
	"github.com/visstt/multisub/backend/internal/config"
	"github.com/visstt/multisub/backend/internal/middleware"
	pkgerrors "github.com/visstt/multisub/backend/pkg/errors"
)

type Server struct {
	app    *fiber.App
	cfg    *config.Config
	db     *pgxpool.Pool
	logger *slog.Logger
}

func New(cfg *config.Config, db *pgxpool.Pool, logger *slog.Logger) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: pkgerrors.Handler,
	})

	s := &Server{
		app:    app,
		cfg:    cfg,
		db:     db,
		logger: logger,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
	s.app.Use(recover.New())
	s.app.Use(requestid.New())
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: []string{s.cfg.CORSOrigins},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))
	s.app.Use(middleware.Logger(s.logger))
}

func (s *Server) setupRoutes() {
	// Health check
	s.app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "multisub-api",
		})
	})

	// API v1
	v1 := s.app.Group("/api/v1")

	// Модуль аутентификации
	authRepo := auth.NewRepository(s.db)
	authSvc := auth.NewService(authRepo, s.cfg)
	authHandler := auth.NewHandler(authSvc)

	authGroup := v1.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.RefreshToken)
	authGroup.Post("/logout", middleware.Auth(s.cfg.JWTSecret), authHandler.Logout)

	// Защищенные маршруты
	protected := v1.Group("", middleware.Auth(s.cfg.JWTSecret))

	// Профиль пользователя
	protected.Get("/me", authHandler.Me)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.cfg.Port)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
