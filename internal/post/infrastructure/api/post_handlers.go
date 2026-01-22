package api

import (
	"encoding/json"
	"net/http"

	"github.com/RLdAB/API-Social-ML/internal/post/application"
	postdomain "github.com/RLdAB/API-Social-ML/internal/post/domain"
)

type PostHandlers struct {
	postService *application.PostService
}

func NewPostHandlers(postService *application.PostService) *PostHandlers {
	return &PostHandlers{postService: postService}
}

// ProductRequest define a estrutura do produto no payload
type ProductRequest struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// CreateProductPost godoc
// @Summary Publicar produto
// @Description Cria uma publicação de produto (não promocional). Apenas sellers podem publicar.
// @Tags Products
// @Accept json
// @Produce json
// @Param body body api.PublishProductRequest true "Payload de publicação"
// @Success 201 {object} api.PostResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 422 {object} api.ErrorResponse
// @Router /products/publish [post]
func (h *PostHandlers) CreateProductPost(w http.ResponseWriter, r *http.Request) {
	var payload postdomain.ProductPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	post, err := h.postService.CreateProduct(payload)
	if err != nil {
		switch err {
		case postdomain.ErrOnlySellerCanPublish:
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		case postdomain.ErrInvalidDate, postdomain.ErrInvalidPayload:
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			// se quiser manter 500 para erro inesperado, troque para StatusInternalServerError
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := toPostResponse(post)
	_ = json.NewEncoder(w).Encode(resp)
}

// CreatePromoProductPost godoc
// @Summary Publicar produto promocional
// @Description Cria uma publicação promocional (has_promo=true e discount > 0). Apenas sellers podem publicar.
// @Tags Promotions
// @Accept json
// @Produce json
// @Param body body api.PublishPromoRequest true "Payload de promoção"
// @Success 201 {object} map[string]string
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 422 {object} api.ErrorResponse
// @Router /products/promo-pub [post]
func (h *PostHandlers) CreatePromoProductPost(w http.ResponseWriter, r *http.Request) {
	var payload postdomain.PromoProductPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := h.postService.CreatePromoProduct(&payload); err != nil {
		// Se você tiver erros tipados no domínio, dá pra mapear 422/404 aqui também.
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "promo published"})
}
