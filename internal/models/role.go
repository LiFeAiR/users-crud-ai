package models

// Role представляет сущность роли
type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	// Связь один ко многим с Permission
	Permissions []Permission `json:"permissions,omitempty"`
}