package application

import (
	"github.com/RLdAB/API-Social-ML/internal/user/domain"
)

// FollowService lida com regras de "following" (US-0001, US-0007)
type FollowService struct {
	repo domain.UserRepository
}

func NewFollowService(repo domain.UserRepository) *FollowService {
	return &FollowService{repo: repo}
}

func (s *FollowService) Execute(userID, sellerID uint) error {
	//Validaçōes obrigatórias (T-0001 - nāo pode seguir a si mesmo)
	if userID == sellerID {
		return domain.ErrSelfFollow
	}

	//Valida existência dos usuários
	if !s.repo.UserExists(userID) || !s.repo.UserExists(sellerID) {
		return domain.ErrUserNotFound
	}

	//Valida se o seller é de fato um vendedor
	seller, err := s.repo.FindByID(sellerID)
	if err != nil {
		return err
	}

	if !seller.IsSeller {
		return domain.ErrNotASeller
	}

	return s.repo.CreateFollow(userID, sellerID)
}

func (s *FollowService) GetFollowersCount(userID uint) (int, error) {
	return s.repo.GetFollowersCount(userID)
}

func (s *FollowService) GetFollowerList(userID uint, order string) ([]domain.User, error) {
	return s.repo.GetFollowerList(userID, order)
}

func (s *FollowService) GetFollowingList(userID uint, order string) ([]domain.User, error) {
	return s.repo.GetFollowingList(userID, order)
}

func (s *FollowService) Unfollow(userID, sellerID uint) error {
	// Valida existência dos usuários
	if !s.repo.UserExists(userID) || !s.repo.UserExists(sellerID) {
		return domain.ErrUserNotFound
	}
	// Valida se o seller é de fato vendedor (opcional)
	seller, err := s.repo.FindByID(sellerID)
	if err != nil {
		return err
	}
	if !seller.IsSeller {
		return domain.ErrNotASeller
	}
	return s.repo.DeleteFollow(userID, sellerID)
}
