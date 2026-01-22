package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID    int     `json:"product_id"`
	Name  string  `json:"product_name"`
	Type  string  `json:"type,omitempty"`
	Brand string  `json:"brand,omitempty"`
	Color string  `json:"color,omitempty"`
	Notes string  `json:"notes,omitempty"`
	Price float64 `json:"price,omitempty"`
}

// Post representa uma publicação promocional
type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID       uint      `json:"user_id" gorm:"index"`
	ProductID    uint      `json:"product_id"`
	ProductName  string    `json:"product_name" gorm:"size:100"`
	ProductType  string    `json:"product_type,omitempty"`
	ProductBrand string    `json:"product_brand,omitempty"`
	Category     string    `json:"category"`
	Content      string    `json:"content" gorm:"type:text"`
	Price        float64   `json:"price" gorm:"type:decimal(10,2)"`
	HasPromo     bool      `json:"has_promo"`
	Discount     float64   `json:"discount" gorm:"type:decimal(5,2)"`
	PromoEndsAt  time.Time `json:"promo_ends_at,omitempty"`
}

type PromoProductPayload struct {
	UserID   uint    `json:"user_id"`
	Product  Product `json:"product"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	HasPromo bool    `json:"has_promo"`
	Discount float64 `json:"discount"`
	Date     string  `json:"date,omitempty"`
}
