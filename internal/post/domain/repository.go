package domain

type ProductPayload struct {
	UserID  uint   `json:"user_id"`
	Date    string `json:"date"`
	Product struct {
		ProductID   int    `json:"product_id"`
		ProductName string `json:"product_name"`
		Type        string `json:"type"`
		Brand       string `json:"brand"`
		Color       string `json:"color"`
		Notes       string `json:"notes"`
	} `json:"product"`
	Category int     `json:"category"`
	Price    float64 `json:"price"`
	HasPromo bool    `json:"has_promo"`
	Discount float64 `json:"discount"`
}

type PostRepository interface {
	// Operaçōes relacionadas a posts
	// Posts (US-0005, US-0006)
	CreatePost(post *Post) error // US-0005
	GetPostsByUser(userID uint) ([]Post, error)
	// Retorna posts recentes dos vendedores que userID segue (últimas N semanas)
	GetRecentFollowedPosts(followerID uint, weeks int, order string) ([]Post, error)
	CountPromotionsBySeller(sellerID uint) (int, error)
	GetRecentPromoPosts(userID uint, weeks int) ([]Post, error)
	GetPromoPostsBySeller(sellerID uint) ([]Post, error)
}
