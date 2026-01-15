package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RLdAB/API-Social-ML/internal/user/application"
	"github.com/RLdAB/API-Social-ML/internal/user/domain"
	"github.com/go-chi/chi/v5"
)

type UserHandlers struct {
	followService *application.FollowService
	userService   *application.UserService
}

// NewUserHandlers é o construtor recomendado para injeçāo de dependências
func NewUserHandlers(fs *application.FollowService, us *application.UserService) *UserHandlers {
	return &UserHandlers{
		followService: fs,
		userService:   us,
	}
}

func (h *UserHandlers) GetFollowersCount(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	count, err := h.followService.GetFollowersCount(userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"followers_count": count,
	})
}

func (h *UserHandlers) GetFollowerList(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	order := r.URL.Query().Get("order")

	list, err := h.followService.GetFollowerList(userID, order)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"followers": list,
	})
}

func (h *UserHandlers) GetFollowingList(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	order := r.URL.Query().Get("order")

	list, err := h.followService.GetFollowingList(userID, order)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"following": list,
	})

}

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
	if err := h.followService.Execute(userID, sellerID); err != nil {
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

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if user.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
	}
	if err := h.userService.CreateUser(&user); err != nil {
		handleServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

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

	if err := h.followService.Unfollow(userID, sellerID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 sem corpo (padrāo REST)
}

func (h *UserHandlers) GetRecentFollowedPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	// Pode aceitar "weeks" como parâmetro query, ex: /products/followed/latest/1?weeks=2
	weeks := 2 // valor padrāo
	if val := r.URL.Query().Get("weeks"); val != "" {
		if wks, err := strconv.Atoi(val); err == nil {
			weeks = wks
		}
	}

	posts, err := h.userService.GetRecentFollowedPosts(userID, weeks)
	if err != nil {
		if err == domain.ErrUserNotFound {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		handleServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts,
	})
}

func (h *UserHandlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post domain.Post
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

func (h *UserHandlers) CountPromotionsBySeller(w http.ResponseWriter, r *http.Request) {
	sellerID, err := strconv.Atoi(chi.URLParam(r, "sellerId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid seller ID")
		return
	}
	count, err := h.userService.CountPromotionsBySeller(sellerID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"seller_id":   sellerID,
		"promo_count": count,
	}
	json.NewEncoder(w).Encode(response)
}

// Faz a chamada ao serviço,
// Retorna um JSON com a lista de usuários.
func (h *UserHandlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers()
	if err != nil {
		handleServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
	})
}

// Recebe o parâmetro userId da URL,
// Converte para int, chama o serviço,
// Retorna erro 404 se não achar, ou 400 se o parâmetro não for número.

func (h *UserHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		handleServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	var updated domain.User
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.userService.UpdateUser(userID, &updated); err != nil {
		if err == domain.ErrUserNotFound {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content para update bem-sucedido
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
