package errors

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

// Sentinel ошибки домена
var (
	ErrNotFound       = errors.New("ресурс не найден")
	ErrConflict       = errors.New("конфликт данных")
	ErrForbidden      = errors.New("доступ запрещён")
	ErrUnauthorized   = errors.New("не авторизован")
	ErrBadRequest     = errors.New("некорректный запрос")
	ErrTooManyReqs    = errors.New("слишком много запросов")
	ErrInternalServer = errors.New("внутренняя ошибка сервера")
)

// AppError — ошибка приложения с HTTP-кодом
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Конструкторы
func NewBadRequest(msg string) *AppError {
	return &AppError{Code: fiber.StatusBadRequest, Message: msg, Err: ErrBadRequest}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{Code: fiber.StatusUnauthorized, Message: msg, Err: ErrUnauthorized}
}

func NewForbidden(msg string) *AppError {
	return &AppError{Code: fiber.StatusForbidden, Message: msg, Err: ErrForbidden}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Code: fiber.StatusNotFound, Message: msg, Err: ErrNotFound}
}

func NewConflict(msg string) *AppError {
	return &AppError{Code: fiber.StatusConflict, Message: msg, Err: ErrConflict}
}

func NewTooManyRequests(msg string) *AppError {
	return &AppError{Code: fiber.StatusTooManyRequests, Message: msg, Err: ErrTooManyReqs}
}

func NewInternal(msg string) *AppError {
	return &AppError{Code: fiber.StatusInternalServerError, Message: msg, Err: ErrInternalServer}
}

// Handler — глобальный обработчик ошибок Fiber
func Handler(c fiber.Ctx, err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.Code).JSON(fiber.Map{
			"error": appErr.Message,
		})
	}

	// Ошибки Fiber (404 и т.д.)
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"error": fiberErr.Message,
		})
	}

	// Непредвиденная ошибка
	slog.Error("необработанная ошибка", "error", err.Error())
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "внутренняя ошибка сервера",
	})
}
