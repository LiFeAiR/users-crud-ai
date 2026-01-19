package models

// Organization represents an organization data structure
type Organization struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	TariffID *int   `json:"tariff_id,omitempty"`
}
