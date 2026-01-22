package api

type PublishProductRequest struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"`

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
}

type PublishPromoRequest struct {
	UserID  uint `json:"user_id"`
	Product struct {
		ProductID   int     `json:"product_id"`
		ProductName string  `json:"product_name"`
		Type        string  `json:"type,omitempty"`
		Brand       string  `json:"brand,omitempty"`
		Color       string  `json:"color,omitempty"`
		Notes       string  `json:"notes,omitempty"`
		Price       float64 `json:"price,omitempty"`
	} `json:"product"`

	Category string  `json:"category"`
	Price    float64 `json:"price"`
	HasPromo bool    `json:"has_promo"`
	Discount float64 `json:"discount"`
	Date     string  `json:"date,omitempty"`
}


