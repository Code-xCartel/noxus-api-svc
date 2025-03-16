package auth

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID int64  `json:"userId"`
	NoxID  string `json:"noxId"`
	jwt.RegisteredClaims
}

type User struct {
	ID        int      `json:"id"`
	NoxID     string   `json:"noxId"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	CreatedAt string   `json:"createdAt"`
	Friends   []string `json:"friends"`
}

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	NoxID    string `json:"noxId"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
