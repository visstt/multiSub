package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/visstt/multisub/backend/internal/config"
	pkgerrors "github.com/visstt/multisub/backend/pkg/errors"
)

// Service — бизнес-логика аутентификации
type Service struct {
	repo Repository
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

// Register создаёт нового пользователя
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// Проверяем, не занят ли email
	existing, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, pkgerrors.NewInternal("ошибка проверки email")
	}
	if existing != nil {
		return nil, pkgerrors.NewConflict("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)
	if err != nil {
		return nil, pkgerrors.NewInternal("ошибка хеширования пароля")
	}

	user := &User{
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Timezone:     "Europe/Moscow",
		Language:     "ru",
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		slog.Error("ошибка создания пользователя", "error", err)
		return nil, pkgerrors.NewInternal("ошибка создания пользователя")
	}

	slog.Info("пользователь зарегистрирован", "user_id", user.ID, "email", user.Email)

	return &RegisterResponse{
		UserID:  user.ID.String(),
		Message: "Регистрация успешна",
	}, nil
}

// Login выполняет вход и возвращает пару токенов
func (s *Service) Login(ctx context.Context, req LoginRequest, clientIP string) (*LoginResponse, error) {
	// Проверка brute-force (5 неудачных попыток за 15 мин)
	since := time.Now().Add(-15 * time.Minute)
	failedCount, err := s.repo.CountRecentFailedAttempts(ctx, req.Email, since)
	if err != nil {
		slog.Error("ошибка подсчёта попыток входа", "error", err)
	}
	if failedCount >= 5 {
		return nil, pkgerrors.NewTooManyRequests("слишком много неудачных попыток входа, попробуйте через 15 минут")
	}

	// Ищем пользователя
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, pkgerrors.NewInternal("ошибка поиска пользователя")
	}
	if user == nil {
		_ = s.repo.RecordLoginAttempt(ctx, req.Email, clientIP, false)
		return nil, pkgerrors.NewUnauthorized("неверный email или пароль")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		_ = s.repo.RecordLoginAttempt(ctx, req.Email, clientIP, false)
		return nil, pkgerrors.NewUnauthorized("неверный email или пароль")
	}

	// Записываем успешную попытку
	_ = s.repo.RecordLoginAttempt(ctx, req.Email, clientIP, true)

	// Генерируем токены
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, pkgerrors.NewInternal("ошибка генерации access token")
	}

	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, pkgerrors.NewInternal("ошибка генерации refresh token")
	}

	slog.Info("пользователь вошёл", "user_id", user.ID, "email", user.Email)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.AccessTokenTTL * 60,
	}, nil
}

// RefreshToken обновляет пару токенов
func (s *Service) RefreshToken(ctx context.Context, refreshTokenRaw string) (*LoginResponse, error) {
	// Хешируем для поиска в БД
	tokenHash := hashToken(refreshTokenRaw)

	storedToken, err := s.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, pkgerrors.NewInternal("ошибка поиска токена")
	}
	if storedToken == nil {
		return nil, pkgerrors.NewUnauthorized("недействительный refresh token")
	}

	// Проверяем, не отозван ли
	if storedToken.RevokedAt != nil {
		return nil, pkgerrors.NewUnauthorized("refresh token отозван")
	}

	// Проверяем срок действия
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, pkgerrors.NewUnauthorized("refresh token истёк")
	}

	// Отзываем старый токен (ротация)
	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		slog.Error("ошибка отзыва старого refresh token", "error", err)
	}

	// Получаем пользователя
	user, err := s.repo.GetUserByID(ctx, storedToken.UserID)
	if err != nil || user == nil {
		return nil, pkgerrors.NewUnauthorized("пользователь не найден")
	}

	// Генерируем новую пару
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, pkgerrors.NewInternal("ошибка генерации access token")
	}

	newRefreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, pkgerrors.NewInternal("ошибка генерации refresh token")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.cfg.AccessTokenTTL * 60,
	}, nil
}

// Logout отзывает все refresh-токены пользователя
func (s *Service) Logout(ctx context.Context, userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return pkgerrors.NewBadRequest("некорректный ID пользователя")
	}
	return s.repo.RevokeAllUserTokens(ctx, uid)
}

// GetMe возвращает данные текущего пользователя
func (s *Service) GetMe(ctx context.Context, userID string) (*MeResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, pkgerrors.NewBadRequest("некорректный ID пользователя")
	}

	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, pkgerrors.NewInternal("ошибка получения профиля")
	}
	if user == nil {
		return nil, pkgerrors.NewNotFound("пользователь не найден")
	}

	return &MeResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}, nil
}

// --- Приватные методы ---

func (s *Service) generateAccessToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  "user",
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Duration(s.cfg.AccessTokenTTL) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *Service) generateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	// Генерируем случайный токен 32 байта
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("генерация случайного токена: %w", err)
	}

	tokenStr := hex.EncodeToString(raw)
	tokenHash := hashToken(tokenStr)

	model := &RefreshTokenModel{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().AddDate(0, 0, s.cfg.RefreshTokenTTL),
	}

	if err := s.repo.SaveRefreshToken(ctx, model); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
