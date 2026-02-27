package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User — модель пользователя в БД
type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Phone        *string    `json:"phone,omitempty"`
	Timezone     string     `json:"timezone"`
	Language     string     `json:"language"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// RefreshTokenModel — модель refresh-токена в БД
type RefreshTokenModel struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// Repository — интерфейс для работы с БД аутентификации
type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	SaveRefreshToken(ctx context.Context, token *RefreshTokenModel) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenModel, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	RecordLoginAttempt(ctx context.Context, email, ip string, success bool) error
	CountRecentFailedAttempts(ctx context.Context, email string, since time.Time) (int, error)
}

// pgRepository — реализация Repository на pgx
type pgRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, timezone, language, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash,
		user.FirstName, user.LastName,
		user.Timezone, user.Language,
		user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("создание пользователя: %w", err)
	}
	return nil
}

func (r *pgRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, timezone, language, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	var u User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash,
		&u.FirstName, &u.LastName, &u.Phone,
		&u.Timezone, &u.Language,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("поиск пользователя по email: %w", err)
	}
	return &u, nil
}

func (r *pgRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, timezone, language, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	var u User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash,
		&u.FirstName, &u.LastName, &u.Phone,
		&u.Timezone, &u.Language,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("поиск пользователя по ID: %w", err)
	}
	return &u, nil
}

func (r *pgRepository) SaveRefreshToken(ctx context.Context, token *RefreshTokenModel) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	token.ID = uuid.New()
	token.CreatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		token.ID, token.UserID, token.TokenHash,
		token.ExpiresAt, token.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("сохранение refresh token: %w", err)
	}
	return nil
}

func (r *pgRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenModel, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	var t RefreshTokenModel
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash,
		&t.ExpiresAt, &t.RevokedAt, &t.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("поиск refresh token: %w", err)
	}
	return &t, nil
}

func (r *pgRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("отзыв refresh token: %w", err)
	}
	return nil
}

func (r *pgRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("отзыв всех токенов пользователя: %w", err)
	}
	return nil
}

func (r *pgRepository) RecordLoginAttempt(ctx context.Context, email, ip string, success bool) error {
	query := `
		INSERT INTO login_attempts (id, email, ip_address, success, created_at)
		VALUES ($1, $2, $3::inet, $4, NOW())
	`
	_, err := r.db.Exec(ctx, query, uuid.New(), email, ip, success)
	if err != nil {
		return fmt.Errorf("запись попытки входа: %w", err)
	}
	return nil
}

func (r *pgRepository) CountRecentFailedAttempts(ctx context.Context, email string, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM login_attempts
		WHERE email = $1 AND success = false AND created_at > $2
	`
	var count int
	err := r.db.QueryRow(ctx, query, email, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("подсчёт неудачных попыток: %w", err)
	}
	return count, nil
}
