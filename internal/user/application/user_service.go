package application

import (
	"strings"
	"time"

	postDomain "github.com/RLdAB/API-Social-ML/internal/post/domain"
	userDomain "github.com/RLdAB/API-Social-ML/internal/user/domain"
)

type UserService struct {
	userRepo userDomain.UserRepository
	postRepo postDomain.PostRepository // Adicionamos a dependência
}

func NewUserService(userRepo userDomain.UserRepository, postRepo postDomain.PostRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (s *UserService) CreatePost(post *postDomain.Post) error {
	// Validações movidas para o PostService
	return s.postRepo.CreatePost(post)
}

func (s *UserService) CreateUser(user *userDomain.User) error {
	// Depois validar outras regras com orientaçāo do Luiz
	if strings.TrimSpace(user.Name) == "" {
		return userDomain.ErrInvalidUser
	}
	return s.userRepo.CreateUser(user)
}

func (s *UserService) GetRecentFollowedPosts(userID uint, weeks int, order string) ([]postDomain.Post, error) {
	// 1. Validaçāo: o usuário existe?
	if !s.userRepo.UserExists(userID) {
		return nil, userDomain.ErrUserNotFound
	}
	// 2. Busca os posts dos sellers seguidos
	return s.postRepo.GetRecentFollowedPosts(userID, weeks, order)
}

func (s *UserService) CountPromotionsBySeller(sellerId uint) (int, error) {
	// Validaçāo: conferir se é vendedor válido (opcional)
	seller, err := s.userRepo.FindByID(sellerId)
	if err != nil {
		return 0, err
	}
	if !seller.IsSeller {
		return 0, userDomain.ErrNotASeller
	}
	return s.postRepo.CountPromotionsBySeller(sellerId)
}

// só repassa a requisição ao repositório
func (s *UserService) ListUsers() ([]userDomain.User, error) {
	return s.userRepo.ListUsers()
}

func (s *UserService) GetUserByID(id uint) (*userDomain.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) UpdateUser(id uint, updated *userDomain.User) error {
	// Validaçōes extra de negócio podem ser colocadas aqui.
	if updated.Name == "" {
		return userDomain.ErrInvalidUser // definir este erro em errors.go
	}
	return s.userRepo.UpdateUser(id, updated)
}

func (s *UserService) CreatePromoProduct(payload *postDomain.PromoProductPayload) error {

	user, err := s.userRepo.FindByID(payload.UserID)
	if err != nil {
		return err
	}
	if !user.IsSeller {
		return userDomain.ErrNotASeller
	}

	if !payload.HasPromo || payload.Discount <= 0 {
		return postDomain.ErrInvalidPromotion
	}
	if payload.Date == "" {
		return postDomain.ErrInvalidDate
	}

	// Parse date
	postDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return postDomain.ErrInvalidDate
	}

	post := &postDomain.Post{
		UserID:      payload.UserID,
		CreatedAt:   postDate,
		ProductName: payload.Product.Name,
		Price:       payload.Price,
		HasPromo:    true,
		Discount:    payload.Discount,
	}

	return s.postRepo.CreatePost(post)
}
