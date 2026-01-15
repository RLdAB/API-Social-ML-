package main

import (
	"github.com/RLdAB/API-Social-ML/internal/user/infrastructure/api"
	"github.com/go-chi/chi/v5"
)

func setupRoutes(r *chi.Mux, userHandlers *api.UserHandlers) {
	//Rotas do User
	r.Post("/users", userHandlers.CreateUser)
	r.Post("/users/{userId}/follow/{sellerId}", userHandlers.FollowUser) // US-0001
	r.Delete("/users/{userId}/follow/{sellerId}", userHandlers.UnfollowUser)
	r.Get("/users/{userId}/followers/list", userHandlers.GetFollowerList)    // US-0003
	r.Get("/users/{userId}/followers/count", userHandlers.GetFollowersCount) // US-0002
	r.Get("/users/{userId}/following/list", userHandlers.GetFollowingList)
	r.Post("/posts", userHandlers.CreatePost)
	r.Get("/products/followed/latest/{userId}", userHandlers.GetRecentFollowedPosts)
	r.Get("/sellers/{sellerId}/promotions/count", userHandlers.CountPromotionsBySeller)
	r.Get("/users", userHandlers.ListUsers)
	r.Get("/users/{userId}", userHandlers.GetUserByID)
	r.Put("/users/{userId}", userHandlers.UpdateUser)
	//r.Get("/products/followed/latest/{userId}", postHandlers.GetRecentPosts) // US-0006
}
