package persistence

import (
	"errors"
	"sort"

	userDomain "github.com/RLdAB/API-Social-ML/internal/user/domain"
	"gorm.io/gorm"
)

// UserRepo implementa a interface UserRepository
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository deve ser exportado (letra maiúscula)
func NewUserRepository(db *gorm.DB) userDomain.UserRepository {
	return &GormUserRepository{db: db}
}

// CreateUser salva um novo usuário
func (r *GormUserRepository) CreateUser(user *userDomain.User) error {
	return r.db.Create(user).Error
}

// FindByID retorna um usuário pelo ID
func (r *GormUserRepository) FindByID(id uint) (*userDomain.User, error) {
	var user userDomain.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// CreateFollow implementa US-0001 (Seguir usuário)
func (r *GormUserRepository) CreateFollow(followerID, sellerID uint) error {
	// Verifica se a relaçāo já existe
	var existing userDomain.Follow
	if err := r.db.Where("follower_id = ? AND seller_id = ?", followerID, sellerID).First(&existing).Error; err == nil {
		return userDomain.ErrAlreadyFollowing
	}

	return r.db.Create(&userDomain.Follow{
		FollowerID: followerID,
		SellerID:   sellerID,
	}).Error
}

// DeleteFollow remove follow entre followerID e sellerID
func (r *GormUserRepository) DeleteFollow(followerID, sellerID uint) error {
	return r.db.Where("follower_id = ? AND seller_id = ?", followerID, sellerID).Delete(&userDomain.Follow{}).Error
}

// GetFollowersList implementa US-0003/0008 (Listar seguidores com ordenaçāo)
func (r *GormUserRepository) GetFollowerList(userID uint, order string) ([]userDomain.User, error) {
	var followers []userDomain.User

	// Query (consulta ao banco de dados) base de join (coluna comum entre duas tabelas) para seguidores
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
		return nil, errors.New("invalid order parameter (use name_asc ou name_desc)")
	}

	return followers, nil
}

// GetFollowingList retorna quem o usuário está seguindo, com ordenaçāo
func (r *GormUserRepository) GetFollowingList(userID uint, order string) ([]userDomain.User, error) {
	var following []userDomain.User
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
		return nil, errors.New("invalid order parameter (use name_asc ou name_desc)")
	}
	return following, nil
}

// GetFollowersCount retorna o número de seguidores de um usuário
func (r *GormUserRepository) GetFollowersCount(userID uint) (int, error) {
	var count int64
	err := r.db.Model(&userDomain.Follow{}).
		Where("seller_id = ?", userID).
		Count(&count).Error
	return int(count), err
}

// UserExists verifica existência do usuário (T-0001)
func (r *GormUserRepository) UserExists(id uint) bool {
	var count int64
	r.db.Model(&userDomain.User{}).
		Where("id = ?", id).
		Count(&count)
	return count > 0
}

func (r *GormUserRepository) UpdateUser(id uint, updated *userDomain.User) error {
	var user userDomain.User
	// Procurar usuário existente
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userDomain.ErrUserNotFound
		}
		return err
	}

	// Atualizar campos (aquilo que pode/merce ser atualizado)
	user.Name = updated.Name
	user.IsSeller = updated.IsSeller

	// Salvar no banco
	return r.db.Save(&user).Error
}

func (r *GormUserRepository) ListUsers() ([]userDomain.User, error) {
	var users []userDomain.User
	// r.db.Find(&users) busca todos os registros da tabela users e preenche o slice
	err := r.db.Find(&users).Error
	return users, err
}
