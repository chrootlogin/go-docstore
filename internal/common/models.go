package common

import "github.com/google/uuid"

const (
	WrongAPIUsageError = "Invalid api call - parameters did not match to method definition"
)

type User struct {
	ID           uuid.UUID `storm:"id"`
	Username     string    `json:"username",storm:"unique"`
	PasswordHash string    `json:"password-hash"`
	Email        string    `json:"email",storm:"unique"`
	IsEnabled    bool      `json:"is-enabled"`
	Permissions  []string  `json:"permissions"`
}

type ApiResponse struct {
	Message string `json:"message"`
}
