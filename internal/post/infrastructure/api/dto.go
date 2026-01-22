package api

import (
	"strconv"
	"time"

	postdomain "github.com/RLdAB/API-Social-ML/internal/post/domain"
	
)

type ProductResponse struct {
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Type        string `json:"type"`
	Brand       string `json:"brand"`
}

type PostResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Date      string    `json:"date"`
	CreatedAt time.Time `json:"created_at"`

	Product  ProductResponse `json:"product"`
	Category int             `json:"category"`
	Price    float64         `json:"price"`
	HasPromo bool            `json:"has_promo"`
	Discount float64         `json:"discount,omitempty"`

	PromoEnds *time.Time `json:"promo_ends_at,omitempty"`
	Content   string     `json:"content,omitempty"`
}

func toPostResponse(p *postdomain.Post) PostResponse {
	cat, _ := strconv.Atoi(p.Category)

	var promoEnds *time.Time
	if !p.PromoEndsAt.IsZero() {
		t := p.PromoEndsAt
		promoEnds = &t
	}

	resp := PostResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		CreatedAt: p.CreatedAt,
		Product: ProductResponse{
			ProductID:   p.ProductID,
			ProductName: p.ProductName,
			Type:        p.ProductType,
			Brand:       p.ProductBrand,
		},
		Category: cat,
		Price:    p.Price,
		HasPromo: p.HasPromo,
		Content:  p.Content,
	}

	if p.HasPromo {
		resp.Discount = p.Discount
		resp.PromoEnds = promoEnds
	}
	return resp
}
