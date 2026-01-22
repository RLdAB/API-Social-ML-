package api

// ErrorResponse padrão de erro
type ErrorResponse struct {
	Error string `json:"error"`
}

type FollowersCountResponse struct {
	FollowersCount int `json:"followers_count"`
}

type SellerPromotionsCountResponse struct {
	SellerID   uint `json:"seller_id"`
	PromoCount int  `json:"promo_count"`
}

type FollowResponse struct {
	Message string `json:"message"`
}

type UsersListResponse struct {
	Users []UserResponse `json:"users"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	IsSeller bool   `json:"is_seller"`
}

type CreatePostRequest struct {
	UserID       uint    `json:"user_id"`
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductType  string  `json:"product_type,omitempty"`
	ProductBrand string  `json:"product_brand,omitempty"`
	Category     string  `json:"category"`
	Price        float64 `json:"price"`
	HasPromo     bool    `json:"has_promo"`
	Discount     float64 `json:"discount"`
	Date         string  `json:"date,omitempty"`
}

type CreatePostResponse struct {
	ID           uint    `json:"id"`
	UserID       uint    `json:"user_id"`
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductType  string  `json:"product_type,omitempty"`
	ProductBrand string  `json:"product_brand,omitempty"`
	Category     string  `json:"category"`
	Content      string  `json:"content"`
	Price        float64 `json:"price"`
	HasPromo     bool    `json:"has_promo"`
	Discount     float64 `json:"discount"`
	CreatedAt    string  `json:"created_at"`
	PromoEndsAt  string  `json:"promo_ends_at,omitempty"`
}

type SimpleUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SimplePost struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user_id"`
	Content   string  `json:"content"`
	Price     float64 `json:"price"`
	HasPromo  bool    `json:"has_promo"`
	Discount  float64 `json:"discount,omitempty"`
	CreatedAt string  `json:"created_at"`
}

type PostResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	Content     string  `json:"content"`
	Price       float64 `json:"price"`
	HasPromo    bool    `json:"has_promo"`
	Discount    float64 `json:"discount,omitempty"`
	CreatedAt   string  `json:"created_at"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
}

type FollowersListResponse struct {
	Followers []UserResponse `json:"followers"`
}

type FollowingListResponse struct {
	Following []UserResponse `json:"following"`
}

type FollowedPostsResponse struct {
	UserID uint               `json:"user_id"`
	Weeks  int                `json:"weeks"`
	Count  int                `json:"count"`
	Posts  []FeedPostResponse `json:"posts"`
}

type FeedPostResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
	Content     string  `json:"content"`
	Price       float64 `json:"price"`
	HasPromo    bool    `json:"has_promo"`
	Discount    float64 `json:"discount,omitempty"`
	CreatedAt   string  `json:"created_at"`
}
