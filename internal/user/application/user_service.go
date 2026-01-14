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
	return s.repo.CreateUser(user)
}
