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
		Joins("JOIN follows ON users.id = follows.seller_id").
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

// GetRecentPromoPosts retorna posts em promoçāo das últimas X semanas
func (r *UserRepo) GetRecentPromoPosts(userID int, weeks int) ([]domain.Promotion, error) {
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

func (r *UserRepo) GetRecentFollowedPosts(userID int, weeks int) ([]domain.Post, error) {
	var posts []domain.Post
	// Pegar a data de corte
	cutoff := time.Now().AddDate(0, 0, -7*weeks)
	// Busque IDs dos vendedores seguidos
	// SELECT p.* FROM posts p
	// JOIN follows f ON p.user_id = f.seller_id
	// WHERE f.follower_id = ? AND p.created_at >- ?

	err := r.db.Table("posts").
		Select("posts.*").
		Joins("JOIN follows ON posts.user_id = follows.seller_id").
		Where("follows.follower_id = ?", userID).
		Where("posts.created_at >= ?", cutoff).
		Order("posts.created_at DESC").
		Find(&posts).Error

	return posts, err
}

func (r *UserRepo) CountPromotionsBySeller(sellerId int) (int, error) {
	var count int64
	err := r.db.Model(&domain.Post{}).
		Where("user_id = ? AND has_promo = ?", sellerId, true).
		Count(&count).Error
	return int(count), err
}

func (r *UserRepo) ListUsers() ([]domain.User, error) {
	var users []domain.User
	// r.db.Find(&users) busca todos os registros da tabela users e preenche o slice
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepo) UpdateUser(id int, updated *domain.User) error {
	var user domain.User
	// Procurar usuário existente
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return err
	}

	// Atualizar campos (aquilo que pode/merce ser atualizado)
	user.Name = updated.Name
	user.IsSeller = updated.IsSeller

	// Salvar no banco
	return r.db.Save(&user).Error
}
