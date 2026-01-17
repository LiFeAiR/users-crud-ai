package models

// Tariff представляет сущность тарифа
type Tariff struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	// Связь один ко многим с Role
	Roles []Role `json:"roles,omitempty"`
}