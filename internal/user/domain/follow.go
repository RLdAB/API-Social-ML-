package domain

import (
	"time"
)

type Follow struct {
	FollowerID int       `gorm:"primaryKey"`
	SellerID   int       `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
