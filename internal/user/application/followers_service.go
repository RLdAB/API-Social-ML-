package application

import "github.com/RLdAB/API-Social-ML/internal/user/domain"

type FollowersService struct {
	repo domain.UserRepository
}

func (s *FollowersService) GetCount(userID int) (int, error) {
	if !s.repo.UserExists(userID) {
		return 0, domain.ErrUserNotFound
	}

	return s.repo.GetFollowersCount(userID)
}
