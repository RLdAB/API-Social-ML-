package persistence

import (
	"time"

	"github.com/RLdAB/API-Social-ML/internal/post/domain"
	postDomain "github.com/RLdAB/API-Social-ML/internal/post/domain"
	"gorm.io/gorm"
)

type GormPostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *GormPostRepository {
	return &GormPostRepository{db: db}
}

// CreatePost salva um novo post
func (r *GormPostRepository) CreatePost(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *GormPostRepository) GetPostsByUser(userID uint) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}

func (r *GormPostRepository) GetRecentFollowedPosts(followerID uint, weeks int, order string) ([]domain.Post, error) {
	var posts []domain.Post
	// Pegar a data de corte
	cutoff := time.Now().AddDate(0, 0, -7*weeks)
	// Busque IDs dos vendedores seguidos
	// SELECT p.* FROM posts p
	// JOIN follows f ON p.user_id = f.seller_id
	// WHERE f.follower_id = ? AND p.created_at >- ?
	query := r.db.Joins("JOIN follows ON posts.user_id = follows.seller_id").
		Where("follows.follower_id = ?", followerID).
		Where("posts.created_at >= ?", cutoff)

	switch order {
	case "date_asc":
		query = query.Order("posts.created_at ASC")
	default: // date_desc
		query = query.Order("posts.created_at DESC")
	}

	err := query.Find(&posts).Error
	return posts, err
}

func (r *GormPostRepository) CountPromotionsBySeller(sellerID uint) (int, error) {
	var count int64
	err := r.db.Model(&domain.Post{}).
		Where("user_id = ? AND has_promo = ?", sellerID, true).
		Count(&count).Error
	return int(count), err
}

// GetRecentPromoPosts retorna posts em promoçāo das últimas X semanas
func (r *GormPostRepository) GetRecentPromoPosts(userID uint, weeks int) ([]domain.Post, error) {
	// Exemplo: buscar promoçōes associadas ao userID nos posts mais recentes (simplicado)
	var posts []postDomain.Post
	// limite de tempo
	cutoff := time.Now().AddDate(0, 0, -7*weeks)
	err := r.db.Where("user_id = ? AND has_promo = ? AND created_at >= ?",
		userID, true, cutoff).Find(&posts).Error

	return posts, err
}
