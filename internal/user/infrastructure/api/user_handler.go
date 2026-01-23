package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	postApp "github.com/RLdAB/API-Social-ML/internal/post/application"
	postDomain "github.com/RLdAB/API-Social-ML/internal/post/domain"
	userApp "github.com/RLdAB/API-Social-ML/internal/user/application"
	"github.com/RLdAB/API-Social-ML/internal/user/domain"
	userDomain "github.com/RLdAB/API-Social-ML/internal/user/domain"
)

type UserHandlers struct {
	followService *userApp.FollowService
	userService   *userApp.UserService
	postService   *postApp.PostService
}

// NewUserHandlers é o construtor recomendado para injeçāo de dependências
func NewUserHandlers(
	fs *userApp.FollowService,
	us *userApp.UserService,
	ps *postApp.PostService) *UserHandlers {
	return &UserHandlers{
		followService: fs,
		userService:   us,
		postService:   ps,
	}
}

// GetFollowersCount godoc
// @Summary Contagem de seguidores
// @Description Retorna a quantidade de seguidores de um vendedor
// @Tags Followers
// @Produce json
// @Param userId path int true "ID do usuário (vendedor)"
// @Success 200 {object} api.FollowersCountResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /users/{userId}/followers/count [get]
func (h *UserHandlers) GetFollowersCount(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	count, err := h.followService.GetFollowersCount(uint(userID))
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"followers_count": count,
	})
}

// GetFollowerList godoc
// @Summary Lista de seguidores
// @Description Retorna a lista de seguidores de um vendedor. Permite ordenação por nome.
// @Tags Followers
// @Produce json
// @Param userId path int true "ID do usuário (vendedor)"
// @Param order query string false "Ordenação" Enums(name_asc,name_desc)
// @Success 200 {object} api.FollowersListResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /users/{userId}/followers/list [get]
func (h *UserHandlers) GetFollowerList(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	order := r.URL.Query().Get("order")

	list, err := h.followService.GetFollowerList(uint(userID), order)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp := make([]UserResponse, 0, len(list))
	for _, u := range list {
		resp = append(resp, toUserResponse(u))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"followers": resp,
	})
}

// GetFollowingList godoc
// @Summary Lista de vendedores seguidos
// @Description Retorna a lista de vendedores seguidos por um usuário. Permite ordenação por nome.
// @Tags Following
// @Produce json
// @Param userId path int true "ID do usuário"
// @Param order query string false "Ordenação" Enums(name_asc,name_desc)
// @Success 200 {object} api.FollowingListResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /users/{userId}/following/list [get]
// @Router /users/{userId}/followed/list [get]
func (h *UserHandlers) GetFollowingList(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	order := r.URL.Query().Get("order")

	list, err := h.followService.GetFollowingList(uint(userID), order)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp := make([]UserResponse, 0, len(list))
	for _, u := range list {
		resp = append(resp, toUserResponse(u))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"following": resp,
	})
}

// FollowUser godoc
// @Summary Seguir vendedor
// @Description Cria relacionamento de follow entre usuário e vendedor
// @Tags Follow
// @Produce json
// @Param userId path int true "ID do usuário (follower)"
// @Param sellerId path int true "ID do vendedor (seller)"
// @Success 201 {object} api.FollowResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 422 {object} api.ErrorResponse
// @Router /users/{userId}/follow/{sellerId} [post]
func (h *UserHandlers) FollowUser(w http.ResponseWriter, r *http.Request) {
	// 1. Parse e validaçāo básica dos IDs
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	sellerID, err := strconv.Atoi(chi.URLParam(r, "sellerId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid seller ID")
		return
	}

	// 2. Chamada do serviço da aplicaçāo
	if err := h.followService.Execute(uint(userID), uint(sellerID)); err != nil {
		handleServiceError(w, err)
		return
	}

	// 3. Resposta de sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "follow relationship created sucessfully",
	})
}

// CreateUser godoc
// @Summary Criar usuário
// @Description Cria um novo usuário (seller ou buyer)
// @Tags Users
// @Accept json
// @Produce json
// @Param body body api.CreateUserRequest true "Dados do usuário"
// @Success 201 {object} api.UserResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /users [post]
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user userDomain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if user.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	if err := h.userService.CreateUser(&user); err != nil {
		if err == userDomain.ErrInvalidUser {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(toUserResponse(user))
}

// UnfollowUser godoc
// @Summary Deixar de seguir vendedor
// @Description Remove relacionamento de follow existente
// @Tags Follow
// @Param userId path int true "ID do usuário (follower)"
// @Param sellerId path int true "ID do vendedor (seller)"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /users/{userId}/follow/{sellerId} [put]
func (h *UserHandlers) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	sellerID, err := strconv.Atoi(chi.URLParam(r, "sellerId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid seller ID")
		return
	}

	if err := h.followService.Unfollow(uint(userID), uint(sellerID)); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 sem corpo (padrāo REST)
}

// GetRecentFollowedPosts godoc
// @Summary Feed de publicações de vendedores seguidos
// @Description Retorna publicações dos vendedores seguidos nas últimas semanas (default=2). Suporta ordenação por data.
// @Tags Feed
// @Produce json
// @Param userId path int true "ID do usuário"
// @Param order query string false "Ordenação" Enums(date_asc,date_desc)
// @Param weeks query int false "Quantidade de semanas (default 2)"
// @Success 200 {object} api.FollowedPostsResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /products/followed/latest/{userId} [get]
// @Router /products/followed/{userId}/list [get]
func (h *UserHandlers) GetRecentFollowedPosts(w http.ResponseWriter, r *http.Request) {
	userID64, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 32)
	if err != nil || userID64 == 0 {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	userID := uint(userID64)

	weeks := 2
	if val := r.URL.Query().Get("weeks"); val != "" {
		wks, err := strconv.Atoi(val)
		if err != nil || wks <= 0 {
			writeError(w, http.StatusBadRequest, "invalid weeks")
			return
		}
		weeks = wks
	}

	order := r.URL.Query().Get("order")

	log.Printf("Fetching posts for user %d, last %d weeks, order %s", userID, weeks, order)

	posts, err := h.userService.GetRecentFollowedPosts(userID, weeks, order)
	if err != nil {
		if err == domain.ErrUserNotFound {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		handleServiceError(w, err)
		return
	}

	respPosts := make([]FeedPostResponse, 0, len(posts))
	for _, p := range posts {
		respPosts = append(respPosts, FeedPostResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			ProductID:   p.ProductID,
			ProductName: p.ProductName,
			Category:    p.Category,
			Content:     p.Content,
			Price:       p.Price,
			HasPromo:    p.HasPromo,
			Discount:    p.Discount,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
		"weeks":   weeks,
		"count":   len(respPosts),
		"posts":   respPosts,
	})
}

// CountPromotionsBySeller godoc
// @Summary Contagem de promoções por vendedor
// @Description Retorna a quantidade de publicações promocionais de um vendedor
// @Tags Promotions
// @Produce json
// @Param sellerId path int true "ID do vendedor"
// @Success 200 {object} api.SellerPromotionsCountResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /sellers/{sellerId}/promotions/count [get]
func (h *UserHandlers) CountPromotionsBySeller(w http.ResponseWriter, r *http.Request) {
	sellerID, err := strconv.ParseUint(chi.URLParam(r, "sellerId"), 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid seller ID")
		return
	}

	count, err := h.userService.CountPromotionsBySeller(uint(sellerID))
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"seller_id":   sellerID,
		"promo_count": count,
	})
}

func (h *UserHandlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post postDomain.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.userService.CreatePost(&post); err != nil {
		handleServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// ListUsers godoc
// @Summary Listar usuários
// @Description Retorna todos os usuários
// @Tags Users
// @Produce json
// @Success 200 {object} api.UsersListResponse
// @Router /users [get]
func (h *UserHandlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers()
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toUserResponse(u))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"users": resp,
	})
}

// GetUserByID godoc
// @Summary Detalhar usuário
// @Description Retorna um usuário pelo ID
// @Tags Users
// @Produce json
// @Param userId path int true "ID do usuário"
// @Success 200 {object} api.UserResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /users/{userId} [get]
func (h *UserHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Convertendo para uint diretamente
	userID, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(uint(userID))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser godoc
// @Summary Atualizar usuário
// @Description Atualiza nome e/ou flag de vendedor
// @Tags Users
// @Accept json
// @Param userId path int true "ID do usuário"
// @Param body body api.CreateUserRequest true "Campos atualizáveis"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /users/{userId} [put]
func (h *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Convertendo para uint
	userID, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var updated domain.User
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.userService.UpdateUser(uint(userID), &updated); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handlerServiceError mapeia erros de domínio para códigos HTTP apropriados
func handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrUserNotFound:
		writeError(w, http.StatusNotFound, err.Error())
	case domain.ErrAlreadyFollowing:
		writeError(w, http.StatusBadRequest, err.Error())
	case domain.ErrSelfFollow:
		writeError(w, http.StatusBadRequest, err.Error())
	case domain.ErrNotASeller:
		writeError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

// writeError é um helper genérico para respostas do erro
func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
