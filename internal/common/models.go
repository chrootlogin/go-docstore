package common

const (
	WrongAPIUsageError = "Invalid api call - parameters did not match to method definition"
)

type User struct {
	ID           int      `storm:"id,increment"` // primary key with auto increment
	Username     string   `json:"username",storm:"unique"`
	PasswordHash string   `json:"password-hash"`
	Email        string   `json:"email",storm:"unique"`
	IsEnabled    bool     `json:"is-enabled"`
	Permissions  []string `json:"permissions"`
}

type ApiResponse struct {
	Message string `json:"message"`
}
