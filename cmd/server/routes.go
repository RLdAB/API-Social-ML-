package main

import (
	"github.com/RLdAB/API-Social-ML/internal/user/infrastructure/api"
	"github.com/go-chi/chi/v5"
)

func setupRoutes(r *chi.Mux, userHandlers *api.UserHandlers) {
	//Rotas do User
	r.Post("/users", userHandlers.CreateUser)                                           // US - 0001
	r.Post("/users/{userId}/follow/{sellerId}", userHandlers.FollowUser)                // US-0001
	r.Get("/users/{userId}/followers/list", userHandlers.GetFollowerList)               // US-0003
	r.Get("/users/{userId}/followers/count", userHandlers.GetFollowersCount)            // US-0002
	r.Get("/users/{userId}/following/list", userHandlers.GetFollowingList)              // US-0004
	r.Post("/posts", userHandlers.CreatePost)                                           // US-0005
	r.Get("/products/followed/latest/{userId}", userHandlers.GetRecentFollowedPosts)    // US-0006
	r.Delete("/users/{userId}/follow/{sellerId}", userHandlers.UnfollowUser)            // US-0007
	r.Get("/users/{userId}/followers/list", userHandlers.GetFollowerList)               // US-0008
	r.Get("/users/{userId}/followed/list", userHandlers.GetFollowingList)               // US-0008
	r.Get("/products/followed/{userId}/list", userHandlers.GetRecentFollowedPosts)      // US-0009
	r.Post("/products/promo-pub", userHandlers.CreatePromoProduct)                      // US-0010
	r.Get("/sellers/{sellerId}/promotions/count", userHandlers.CountPromotionsBySeller) // US-0011
	r.Get("/users", userHandlers.ListUsers)
	r.Get("/users/{userId}", userHandlers.GetUserByID)
	r.Put("/users/{userId}", userHandlers.UpdateUser)
	//r.Get("/products/followed/latest/{userId}", postHandlers.GetRecentPosts) // US-0006
}
