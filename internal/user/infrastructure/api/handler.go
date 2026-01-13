package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RLdAB/API-Social-ML/internal/user/application"
	"github.com/go-chi/chi/v5"
)

type UserHandlers struct {
	followService *application.FollowService
}

// NewUserHandlers é o construtor recomendado para injeçāo de dependências
func NewUserHandlers(fs *application.FollowService) *UserHandlers {
	return &UserHandlers{
		followService: fs,
	}
}

func (h *UserHandlers) GetFollowersCount(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
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
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
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
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	order := r.URL.Query().Get("order")

	list, err := s.followService.GetFollowingList(userID, order)
	if err != nil {
		handleServiceError(w, err)
		return
	}
  
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface(){
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

// handlerServiceError mapeia erros de domínio para códigos HTTP apropriados
func handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case application.ErrUserNotFound:
		writeError(w, http.StatusNotFound, err.Error())
	case application.ErrSelfFollow:
		writeError(w, http.StatusBadRequest, err.Error())
	case application.ErrNotASeller:
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
