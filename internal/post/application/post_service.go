package application

import (
	"errors"
	"fmt"
	"strings"
	"time"

	postdomain "github.com/RLdAB/API-Social-ML/internal/post/domain"
	userDomain "github.com/RLdAB/API-Social-ML/internal/user/domain"
)

type PostService struct {
	postRepo postdomain.PostRepository
	userRepo userDomain.UserRepository
}

func NewPostService(
	postRepo postdomain.PostRepository,
	userRepo userDomain.UserRepository,
) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

type ProductInput struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (s *PostService) CreatePost(post *postdomain.Post) error {
	// 1. Validação do usuário
	user, err := s.userRepo.FindByID(post.UserID)
	if err != nil {
		return err
	}

	if !user.IsSeller {
		return errors.New("only sellers can create posts")
	}

	// 2. Validação do post
	if post.ProductName == "" || post.Price <= 0 {
		return errors.New("product name and price are required")
	}

	// 3. Validação de promoção
	if post.HasPromo {
		if post.Discount <= 0 || post.Discount >= 1 {
			return errors.New("invalid discount value")
		}
	} else {
		post.Discount = 0
	}

	// 4. Definir timestamps
	if post.CreatedAt.IsZero() {
		post.CreatedAt = time.Now()
	}

	return s.postRepo.CreatePost(post)
}

func (s *PostService) CreatePromoProduct(payload *postdomain.PromoProductPayload) error {
	userID := payload.UserID

	post := postdomain.Post{
		UserID:       userID,
		ProductName:  payload.Product.Name,
		ProductType:  payload.Product.Type,
		ProductBrand: payload.Product.Brand,
		Category:     payload.Category,
		Price:        payload.Price,
		HasPromo:     true,
		Discount:     payload.Discount,
		CreatedAt:    time.Now(),
	}

	if payload.Date != "" {
		if date, err := time.Parse("2006-01-02", payload.Date); err == nil {
			post.CreatedAt = date
		}
	}

	if payload.Discount <= 0 || payload.Discount >= 1 {
		return errors.New("discount must be between 0 and 1")
	}

	return s.CreatePost(&post)
}

func (s *PostService) CreateProduct(p postdomain.ProductPayload) (*postdomain.Post, error) {
	// valida user_id
	if p.UserID <= 0 {
		return nil, postdomain.ErrInvalidPayload
	}
	uid := uint(p.UserID)

	// user existe?
	user, err := s.userRepo.FindByID(uid)
	if err != nil {
		return nil, err // aqui pode virar 404 ou 500 dependendo do repo
	}

	// só seller publica
	if !user.IsSeller {
		return nil, postdomain.ErrOnlySellerCanPublish
	}

	// validaçōes básicas
	if strings.TrimSpace(p.Product.ProductName) == "" || p.Price <= 0 {
		return nil, postdomain.ErrInvalidPayload
	}

	post := postdomain.Post{
		UserID:       uid,
		ProductID:    uint(p.Product.ProductID),
		ProductName:  p.Product.ProductName,
		ProductType:  p.Product.Type,
		ProductBrand: p.Product.Brand,
		Category:     fmt.Sprintf("%d", p.Category), // converte int -> string
		Price:        p.Price,
		HasPromo:     false,
		Discount:     0,
		Content:      fmt.Sprintf("Publicaçāo: %s", p.Product.ProductName),
		CreatedAt:    time.Now(),
	}

	// data no formato dd-MM-aaaa
	if p.Date != "" {
		d, err := time.Parse("02-01-2006", p.Date)
		if err != nil {
			return nil, postdomain.ErrInvalidDate
		}
		post.CreatedAt = d
	}

	if err := s.CreatePost(&post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) ListPromoPostsBySeller(userID uint) ([]postdomain.Post, string, error) {
	// 1) validar user existe e é seller
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, "", err
	}
	if !user.IsSeller {
		return nil, "", errors.New("user is not a seller")
	}

	// 2) buscar promo posts no repo
	posts, err := s.postRepo.GetPromoPostsBySeller(userID)
	if err != nil {
		return nil, "", err
	}

	return posts, user.Name, nil
}
