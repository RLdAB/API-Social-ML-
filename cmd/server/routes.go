package main

import (
	"github.com/RLdAB/API-Social-ML/internal/user/infrastructure/api"
	"github.com/go-chi/chi/v5"
)

func setupRoutes(r *chi.Mux, userHandlers *api.UserHandlers) {
	//Rotas do User
	r.Post("/users/{userId}/follow/{sellerId}", userHandlers.FollowUser)     // US-0001
	r.Get("/users/{userId}/followers/count", userHandlers.GetFollowersCount) // US-0002
	r.Get("/users/{userId}/followers/list", userHandlers.GetFollowerList)    // US-0003
	r.Post("/users", userHandlers.CreateUser)
	//r.Get("/products/followed/latest/{userId}", postHandlers.GetRecentPosts) // US-0006
}
