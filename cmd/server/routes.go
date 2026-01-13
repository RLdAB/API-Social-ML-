package main

import (
	"github.com/RLdAB/API-Social-ML/internal/user/application"
	"github.com/RLdAB/API-Social-ML/internal/user/infrastructure/api"
	"github.com/RLdAB/API-Social-ML/internal/user/infrastructure/persistence"
	"github.com/go-chi/chi/v5"
)

func setupRoutes() *chi.Mux {
	r := chi.NewRouter()

	//Injeçāo de dependências
	userRepo := persistence.NewUserRepository() //Implementar isto
	followService := application.NewFollowService(userRepo)
	userHandlers := api.NewUserHandlers(followService)

	//Rotas do User
	r.Post("/users/{userId}/follow/{sellerId}", userHandlers.FollowUser)     // US-0001
	r.Get("/users/{userId}/followers/count", userHandlers.GetFollowersCount) // US-0002
	r.Get("users/{userId}/followers/list", userHandlers.GetFollowersList)    // US-0003

	//Post routes
	r.Get("/products/followed/latest/{userId}", postHandlers.GetRecentPosts) // US-0006
	return r
}
