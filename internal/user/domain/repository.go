package domain

import (
	"time"
)

//Repository.go - Interface completa para operaçōes de persistência

type UserRepository interface {
	//Operaçōes básicas de usuário
	CreateUser(user *User) error
	FindByID(id int) (*User, error)
	UserExists(id int) bool

	//Seguidores/Seguindo (US-0001 a US-0004, US-0007)
	CreateFollow(followerID, sellerID int) error               // US-0001
	DeleteFollow(followerID, sellerID int) error               // US-0007
	GetFollowersCount(userID int) (int, error)                 // USS-0002
	GetFollowerList(userID int, order string) ([]User, error)  //US-0003, US-0008
	GetFollowingList(userID int, order string) ([]User, error) // US-0004

	// Posts (US-0005, US-0006)
	CreatePost(post *Post) error // US-0005
	GetRecentPosts(userID int, weeks int) ([]Promotion, error)
}

// Estruturas auxiliares

type Post struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int       `gorm:"index"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	//Campos para US-0006
	Category int `gorm:"index"`
	Price    float64
	HasPromo bool
	Discount float64
}

type Promotion struct {
	PostID    int `gorm:"primaryKey"`
	ExpiresAt time.Time
}
