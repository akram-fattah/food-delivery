package models

import "time"

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	IsVerified bool `json:"is_verified"`
	VerificationCode string    `json:"verification_code"`
	CreateAt time.Time `json:"created_at"`
}

type UpdateUserRequest struct {
	Name             *string `json:"name"`
	Username         *string `json:"username"`
	Email            *string `json:"email"`
	Password         *string `json:"password"`
}

type UpdateRoleRequest struct {
	Role *string `json:"role"`
}