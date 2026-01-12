package domain

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" validate:"required,max=15"`
	IsSeller bool   `json:"is_seller"`
	//Adicionar outros campos depois do check do Luiz
}
