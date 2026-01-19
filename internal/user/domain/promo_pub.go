package domain

type PromoProductPayload struct {
	UserID  int    `json:"user_id"`
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
