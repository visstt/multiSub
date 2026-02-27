package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/visstt/multisub/backend/internal/config"
	"github.com/visstt/multisub/backend/internal/database"
	"github.com/visstt/multisub/backend/internal/server"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		slog.Error("не удалось загрузить конфигурацию", "error", err)
		os.Exit(1)
	}

	// Настройка логгера
	logger := setupLogger(cfg.Env)
	slog.SetDefault(logger)

	slog.Info("запуск MultiSub API", "env", cfg.Env, "port", cfg.Port)

	// Подключение к БД
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("не удалось подключиться к БД", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("подключение к PostgreSQL установлено")

	// Запуск HTTP-сервера
	srv := server.New(cfg, db, logger)

	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("ошибка сервера", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("завершение работы сервера...")
	if err := srv.Shutdown(); err != nil {
		slog.Error("ошибка при остановке сервера", "error", err)
	}
	slog.Info("сервер остановлен")
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case "production":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}
