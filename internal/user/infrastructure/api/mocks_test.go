package api

import (
	postDomain "github.com/RLdAB/API-Social-ML/internal/post/domain"
	userDomain "github.com/RLdAB/API-Social-ML/internal/user/domain"
)

type mockUserRepo struct {
	created   []*userDomain.User
	createErr error
}

func (m *mockUserRepo) CreateUser(user *userDomain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.created = append(m.created, user)
	return nil
}

func (m *mockUserRepo) FindByID(id uint) (*userDomain.User, error)   { return nil, nil }
func (m *mockUserRepo) UserExists(id uint) bool                      { return true }
func (m *mockUserRepo) CreateFollow(followerID, sellerID uint) error { return nil }
func (m *mockUserRepo) DeleteFollow(followerID, sellerID uint) error { return nil }
func (m *mockUserRepo) GetFollowersCount(userID uint) (int, error)   { return 0, nil }
func (m *mockUserRepo) GetFollowerList(userID uint, order string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetFollowingList(userID uint, order string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) ListUsers() ([]userDomain.User, error)           { return nil, nil }
func (m *mockUserRepo) UpdateUser(id uint, user *userDomain.User) error { return nil }

type mockPostRepo struct{}

func (m *mockPostRepo) CreatePost(post *postDomain.Post) error                { return nil }
func (m *mockPostRepo) GetPostsByUser(userID uint) ([]postDomain.Post, error) { return nil, nil }
func (m *mockPostRepo) GetRecentFollowedPosts(followerID uint, weeks int, order string) ([]postDomain.Post, error) {
	return nil, nil
}
func (m *mockPostRepo) CountPromotionsBySeller(sellerID uint) (int, error) { return 0, nil }
func (m *mockPostRepo) GetRecentPromoPosts(userID uint, weeks int) ([]postDomain.Post, error) {
	return nil, nil
}
func (m *mockPostRepo) GetPromoPostsBySeller(sellerID uint) ([]postDomain.Post, error) {
	return nil, nil
}
