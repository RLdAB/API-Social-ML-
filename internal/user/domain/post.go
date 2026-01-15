package domain

import (
	"time"
)

type Post struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"index"` // o vendedor que faz o post
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Category  int       `json:"category" gorm:"index"`
	Price     float64   `json:"price"`
	HasPromo  bool      `json:"has_promo"`
	Discount  float64   `json:"discount"`
}

type Promotion struct {
	PostID    int `gorm:"primaryKey"`
	ExpiresAt time.Time
}
