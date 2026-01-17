package models

// User represents a user data structure
type User struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Email          string        `json:"email"`
	PasswordHash   string        `json:"password_hash"`
	Organization   *Organization `json:"organization"`
}
