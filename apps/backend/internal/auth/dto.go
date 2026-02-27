package auth

import "time"

// --- Request DTOs ---

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *RegisterRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if r.Email == "" {
		errs["email"] = "email обязателен"
	}
	if len(r.Password) < 8 {
		errs["password"] = "пароль должен содержать минимум 8 символов"
	}
	if len(r.Password) > 72 {
		errs["password"] = "пароль не может быть длиннее 72 символов"
	}
	if r.FirstName == "" {
		errs["first_name"] = "имя обязательно"
	}
	if r.LastName == "" {
		errs["last_name"] = "фамилия обязательна"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if r.Email == "" {
		errs["email"] = "email обязателен"
	}
	if r.Password == "" {
		errs["password"] = "пароль обязателен"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// --- Response DTOs ---

type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // секунды до истечения access token
}

type MeResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}
