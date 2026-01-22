package domain

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time      `json:"-"`
	UpdateddAt time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	Name     string `json:"name" validate:"required,max=15"` // Validaçāo da documentaçāo
	IsSeller bool   `json:"is_seller"`
	//Adicionar outros campos depois do check do Luiz
}

type Follow struct {
	FollowerID uint `gorm:"primaryKey;autoIncrement:false"`
	SellerID   uint `gorm:"primaryKey;autoIncrement:false"`
}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrAlreadyFollowing = errors.New("already following")
)

// Validaçāo manual para a US-0001
func (u *User) CanFollow(target *User) error {
	if u.ID == target.ID {
		return errors.New("cannot follow yourself")
	}
	if !target.IsSeller {
		return errors.New("target is not a seller")
	}
	return nil
}
