package models

// User represents a user data structure
type User struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	Password     string        `json:"password,omitempty"`
	Organization *Organization `json:"organization"`
}
