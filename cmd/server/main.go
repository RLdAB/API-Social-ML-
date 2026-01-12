package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RLdAB/API-Social-ML/internal/user/infrastructure/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	//Configurar rotas (ex:user)
	userHandler := api.NewUserHandler( /* dependencies */ )
	r.Post("/users/{userId}/follow/{userIdToFollow}", userHandler.FollowUser)

	//Iniciar servidor
	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
