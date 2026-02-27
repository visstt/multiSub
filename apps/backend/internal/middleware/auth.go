package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	pkgerrors "github.com/visstt/multisub/backend/pkg/errors"
)

// Auth — middleware для проверки JWT access token
func Auth(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return pkgerrors.NewUnauthorized("токен авторизации отсутствует")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return pkgerrors.NewUnauthorized("неверный формат токена (ожидается Bearer <token>)")
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, pkgerrors.NewUnauthorized("неподдерживаемый алгоритм подписи")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			return pkgerrors.NewUnauthorized("невалидный или просроченный токен")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return pkgerrors.NewUnauthorized("невалидные данные токена")
		}

		// Сохраняем данные пользователя в контексте
		c.Locals("userID", claims["sub"])
		c.Locals("email", claims["email"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
