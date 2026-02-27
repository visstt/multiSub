package auth

import (
	"github.com/gofiber/fiber/v3"

	pkgerrors "github.com/visstt/multisub/backend/pkg/errors"
)

// Handler — HTTP-обработчики аутентификации
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Register godoc
// POST /api/v1/auth/register
func (h *Handler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return pkgerrors.NewBadRequest("некорректный JSON")
	}

	if errs := req.Validate(); errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":  "ошибка валидации",
			"fields": errs,
		})
	}

	resp, err := h.svc.Register(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// Login godoc
// POST /api/v1/auth/login
func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return pkgerrors.NewBadRequest("некорректный JSON")
	}

	if errs := req.Validate(); errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":  "ошибка валидации",
			"fields": errs,
		})
	}

	resp, err := h.svc.Login(c.Context(), req, c.IP())
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

// RefreshToken godoc
// POST /api/v1/auth/refresh
func (h *Handler) RefreshToken(c fiber.Ctx) error {
	var req RefreshRequest
	if err := c.Bind().JSON(&req); err != nil {
		return pkgerrors.NewBadRequest("некорректный JSON")
	}

	if req.RefreshToken == "" {
		return pkgerrors.NewBadRequest("refresh_token обязателен")
	}

	resp, err := h.svc.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

// Logout godoc
// POST /api/v1/auth/logout
func (h *Handler) Logout(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return pkgerrors.NewUnauthorized("невалидный токен")
	}

	if err := h.svc.Logout(c.Context(), userID); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "вы вышли из системы"})
}

// Me godoc
// GET /api/v1/me
func (h *Handler) Me(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return pkgerrors.NewUnauthorized("невалидный токен")
	}

	resp, err := h.svc.GetMe(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}
