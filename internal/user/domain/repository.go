package domain

//Repository.go - Interface completa para operaçōes de persistência

type UserRepository interface {
	//Operaçōes básicas de usuário
	CreateUser(user *User) error
	// Retorna um usuário via ID, ou erro se não encontrado.
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
	GetRecentPromoPosts(userID int, weeks int) ([]Promotion, error)
	// Retorna posts recentes dos vendedores que userID segue (últimas N semanas)
	GetRecentFollowedPosts(userID int, weeks int, order string) ([]Post, error)
	CountPromotionsBySeller(sellerId int) (int, error)
	ListUsers() ([]User, error)
	UpdateUser(id int, user *User) error
}
