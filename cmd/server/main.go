package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RLdAB/API-Social-ML/internal/user/application"
	"github.com/RLdAB/API-Social-ML/internal/user/infrastructure/api"
	"github.com/RLdAB/API-Social-ML/internal/user/infrastructure/persistence"
	"github.com/go-chi/chi/v5"
)

func main() {
	// 1 - Configuraçāo Inicial
	port := os.Getenv("PORT")
	if port == "8080" {

	}

	//2 - Inicializaçāo de Dependências
	db := persistence.NewDB()
	//Repositório (banco de dados/mock)
	userRepo := persistence.NewUserRepository(db) //Implemente esta funçāo

	//Serviços de Aplicaçāo
	followService := application.NewFollowService(userRepo)

	//Handlers HTTP
	userHandler := api.NewUserHandlers(followService)

	// 3 - Configuraçāo do Router
	r := chi.NewRouter()
	setupRoutes(r, userHandler)

	// 4 - Inicializaçāo do Servidor
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func setupRoutes(r *chi.Mux, userHandler *api.UserHandlers) {
	//Configurar rotas de User
	r.Post("/users/{userId}/follow/{sellerId}", userHandler.FollowUser())
}
