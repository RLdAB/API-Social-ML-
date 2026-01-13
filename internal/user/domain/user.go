package domain

import "errors"

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" validate:"required,max=15"` // Validaçāo da documentaçāo
	IsSeller bool   `json:"is_seller"`
	//Adicionar outros campos depois do check do Luiz
}

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
