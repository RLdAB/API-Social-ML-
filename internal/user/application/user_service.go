package application

import (
	"github.com/RLdAB/API-Social-ML/internal/user/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *domain.User) error {
	// Depois validar outras regras com orientaçāo do Luiz
	if user.Name == "" {
		return domain.ErrInvalidUser
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) CreatePost(post *domain.Post) error {
	// 1. Validaçāo
	if post.Content == "" || post.UserID == 0 {
		return domain.ErrInvalidPost // criar esse erro no domain/errors.go
	}
	// 2. Checa se o usuário existe e é vendedor
	user, err := s.repo.FindByID(post.UserID)
	if err != nil {
		return err
	}
	if user == nil || !user.IsSeller {
		return domain.ErrNotASeller // criar esse erro com errors.go
	}
	// 3. Salvar
	return s.repo.CreatePost(post)
}

func (s *UserService) GetRecentFollowedPosts(userID int, weeks int, order string) ([]domain.Post, error) {
	// 1. Validaçāo: o usuário existe?
	if !s.repo.UserExists(userID) {
		return nil, domain.ErrUserNotFound
	}
	// 2. Busca os posts dos sellers seguidos
	return s.repo.GetRecentFollowedPosts(userID, weeks, order)
}

func (s *UserService) CountPromotionsBySeller(sellerId int) (int, error) {
	// Validaçāo: conferir se é vendedor válido (opcional)
	seller, err := s.repo.FindByID(sellerId)
	if err != nil {
		return 0, err
	}
	if !seller.IsSeller {
		return 0, domain.ErrNotASeller
	}
	return s.repo.CountPromotionsBySeller(sellerId)
}

// só repassa a requisição ao repositório
func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.repo.ListUsers()
}

func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) UpdateUser(id int, updated *domain.User) error {
	// Validaçōes extra de negócio podem ser colocadas aqui.
	if updated.Name == "" {
		return domain.ErrInvalidUser // definir este erro em errors.go
	}
	return s.repo.UpdateUser(id, updated)
}
