package main

import (
	postapi "github.com/RLdAB/API-Social-ML/internal/post/infrastructure/api"
	userapi "github.com/RLdAB/API-Social-ML/internal/user/infrastructure/api"
	"github.com/go-chi/chi/v5"
)

func setupRoutes(r *chi.Mux, userHandlers *userapi.UserHandlers, postHandlers *postapi.PostHandlers) {
	// ========== USERS ==========
	r.Route("/users", func(r chi.Router) {
		// CRUD
		r.Post("/", userHandlers.CreateUser)
		r.Get("/", userHandlers.ListUsers)
		r.Get("/{userId}", userHandlers.GetUserByID)
		r.Put("/{userId}", userHandlers.UpdateUser)

		// Follow / Unfollow (US-0001 / US-0007)
		r.Post("/{userId}/follow/{sellerId}", userHandlers.FollowUser)
		r.Put("/{userId}/follow/{sellerId}", userHandlers.UnfollowUser)

		// Followers (US-0002 / US-0003 / US-0008 via ?order=name_asc|name_desc)
		r.Get("/{userId}/followers/count", userHandlers.GetFollowersCount)
		r.Get("/{userId}/followers/list", userHandlers.GetFollowerList)

		// Following (US-0004 / US-0008 via ?order=name_asc|name_desc)
		r.Get("/{userId}/following/list", userHandlers.GetFollowingList)
		// Alias compatível com enunciado
		r.Get("/{userId}/followed/list", userHandlers.GetFollowingList)
	})

	// Feed de seguidos (US-0006 / US-0009)
	r.Route("/products/followed", func(r chi.Router) {
		r.Get("/latest/{userId}", userHandlers.GetRecentFollowedPosts) // latest (últimas 2 semanas)
		r.Get("/{userId}/list", userHandlers.GetRecentFollowedPosts)   // list + order=date_asc|date_desc
	})

	// Produtos (US-0005) e Promoções (US-0010)
	r.Route("/products", func(r chi.Router) {
		r.Post("/publish", postHandlers.CreateProductPost)
		r.Post("/promo-pub", postHandlers.CreatePromoProductPost)
	})

	// Métricas promo (US-0011)
	r.Get("/sellers/{sellerId}/promotions/count", userHandlers.CountPromotionsBySeller)
}
