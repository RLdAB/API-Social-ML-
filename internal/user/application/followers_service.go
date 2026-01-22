package application

import "github.com/RLdAB/API-Social-ML/internal/user/domain"

type FollowersService struct {
	repo domain.UserRepository
}

func (s *FollowersService) GetCount(userID uint) (int, error) {
	if !s.repo.UserExists(userID) {
		return 0, domain.ErrUserNotFound
	}

	return s.repo.GetFollowersCount(userID)
}

func (s *FollowersService) GetFollowingList(userID uint, order string) ([]domain.User, error) {
	// Opcional - Adicionei validaçāo se o usuário existe
	if !s.repo.UserExists(userID) {
		return nil, domain.ErrUserNotFound
	}
	return s.repo.GetFollowingList(userID, order)
}


