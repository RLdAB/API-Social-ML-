package persistence

import (
	"errors"
	"sort"
	"time"

	"github.com/RLdAB/API-Social-ML/internal/user/domain"
	"gorm.io/gorm"
)

// UserRepo implementa a interface UserRepository
type UserRepo struct {
	db *gorm.DB
}

// GetFollowerList implements [domain.UserRepository].
func (r *UserRepo) GetFollowerList(userID int, order string) ([]domain.User, error) {
	panic("unimplemented")
}

// GetRecentPosts implements [domain.UserRepository].
func (r *UserRepo) GetRecentPosts(userID int, weeks int) ([]domain.Promotion, error) {
	panic("unimplemented")
}

// NewUserRepository é o construtor (inicializa o repositório)
func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// CreateUser salva um novo usuário
func (r *UserRepo) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByID retorna um usuário pelo ID
func (r *UserRepo) FindByID(id int) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// CreateFollow implementa US-0001 (Seguir usuário)
func (r *UserRepo) CreateFollow(followerID, sellerID int) error {
	// Verifica se a relaçāo já existe
	var existing domain.Follow
	if err := r.db.Where("follower_id = ? AND seller_id = ?", followerID, sellerID).First(&existing).Error; err == nil {
		return domain.ErrAlreadyFollowing
	}

	return r.db.Create(&domain.Follow{
		FollowerID: followerID,
		SellerID:   sellerID,
	}).Error
}

// DeleteFollow remove follow entre followerID e sellerID
func (r *UserRepo) DeleteFollow(followerID, sellerID int) error {
	return r.db.Where("follower_id = ? AND seller_id = ?", followerID, sellerID).Delete(&domain.Follow{}).Error
}

// GetFollowersList implementa US-0003/0008 (Listar seguidores com ordenaçāo)
func (r *UserRepo) GetFollowerList(userID int, order string) ([]domain.User, error) {
	var followers []domain.User

	// Query base de join para seguidores
	query := r.db.Table("users").
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.seller_id = ?", userID)

	// Executa a query
	if err := query.Find(&followers).Error; err != nil {
		return nil, err
	}

	// Ordenaçāo (US-0008)
	switch order {
	case "name_asc":
		sort.Slice(followers, func(i, j int) bool {
			return followers[i].Name < followers[j].Name
		})
	case "name_desc":
		sort.Slice(followers, func(i, j int) bool {
			return followers[i].Name > followers[j].Name
		})
	default:
		// Sem ordenaçāo ou ordenaçāo inválida, só retorna
		return nil, errors.New("invalid order parameter")
	}

	return followers, nil
}

// GetFollowingList retorna quem o usuário está seguindo, com ordenaçāo
func (r *UserRepo) GetFollowingList(userID int, order string) ([]domain.User, error) {
	var following []domain.User
	// join para seguidos
	query := r.db.Table("users").
		Joins("JOIN follows ON users.id = follows.seller.id").
		Where("follows.follower_id = ?", userID)
	if err := query.Find(&following).Error; err != nil {
		return nil, err
	}
	switch order {
	case "name_asc":
		sort.Slice(following, func(i, j int) bool {
			return following[i].Name < following[j].Name
		})
	case "name_desc":
		sort.Slice(following, func(i, j int) bool {
			return following[i].Name > following[j].Name
		})
	default:
		//Sem ordenaçāo ou ordenaçāo inválida, só retorna
	}
	return following, nil
}

// GetFollowersCount retorna o número de seguidores de um usuário
func (r *UserRepo) GetFollowersCount(userID int) (int, error) {
	var count int64
	err := r.db.Model(&domain.Follow{}).
		Where("seller_id = ?", userID).
		Count(&count).Error
	return int(count), err
}

// UserExists verifica existência do usuário (T-0001)
func (r *UserRepo) UserExists(id int) bool {
	var count int64
	r.db.Model(&domain.User{}).
		Where("id = ?", id).
		Count(&count)
	return count > 0
}

// CreatePost salva um novo post
func (r *UserRepo) CreatePost(post *domain.Post) error {
	return r.db.Create(post).Error
}

// GetRecentPosts retorna posts em promoçāo das últimas X semanas
func (r *UserRepo) GetRecentePosts(userID int, weeks int) ([]domain.Promotion, error) {
	// Exemplo: buscar promoçōes associadas ao userID nos posts mais recentes (simplicado)
	var promotions []domain.Promotion
	// limite de tempo
	cutoff := time.Now().AddDate(0, 0, -7*weeks)
	err := r.db.Table("promotions").
		Joins("JOIN posts ON promotions.post_id = posts.id").
		Where("posts.user_id = ? AND posts.created_at >= ?", userID, cutoff).
		Find(&promotions).Error
	return promotions, err
}
