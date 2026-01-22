package domain

//Repository.go - Interface completa para operaçōes de persistência

type UserRepository interface {
	//Operaçōes básicas de usuário
	CreateUser(user *User) error
	// Retorna um usuário via ID, ou erro se não encontrado.
	FindByID(id uint) (*User, error)
	UserExists(id uint) bool

	//Seguidores/Seguindo (US-0001 a US-0004, US-0007)
	CreateFollow(followerID, sellerID uint) error               // US-0001
	DeleteFollow(followerID, sellerID uint) error               // US-0007
	GetFollowersCount(userID uint) (int, error)                 // USS-0002
	GetFollowerList(userID uint, order string) ([]User, error)  //US-0003, US-0008
	GetFollowingList(userID uint, order string) ([]User, error) // US-0004

	// Operaçōes adicionais
	ListUsers() ([]User, error)
	UpdateUser(id uint, user *User) error
}
