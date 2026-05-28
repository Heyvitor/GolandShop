package model

import "time"

const (
	RoleAdmin  = "admin"
	RoleUser   = "user"   // Dono de loja
	RoleClient = "client" // Consumidor
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Store struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"-"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Item struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id,omitempty"`
	StoreID      string    `json:"store_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	Variant      string    `json:"variant"`
	VariantPrice float64   `json:"variant_price"`
	ShippingType string    `json:"shipping_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Order struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	StoreID     string    `json:"store_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
