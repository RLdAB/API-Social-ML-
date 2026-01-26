package application

import postDomain "github.com/RLdAB/API-Social-ML/internal/post/domain"

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
