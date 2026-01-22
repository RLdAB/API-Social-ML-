package main

import (
	"log"
	"net/http"
	"os"

	postApplication "github.com/RLdAB/API-Social-ML/internal/post/application"
	postApi "github.com/RLdAB/API-Social-ML/internal/post/infrastructure/api"
	"github.com/RLdAB/API-Social-ML/internal/post/infrastructure/persistence"
	userApplication "github.com/RLdAB/API-Social-ML/internal/user/application"
	userApi "github.com/RLdAB/API-Social-ML/internal/user/infrastructure/api"
	userPersistence "github.com/RLdAB/API-Social-ML/internal/user/infrastructure/persistence"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1 - Configuraçāo Inicial
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//2 - Inicializaçāo de Dependências
	log.Println("Porta definida")
	db := userPersistence.NewDB()
	log.Println("Banco conectado")
	//Repositório (banco de dados/mock)
	userRepo := userPersistence.NewUserRepository(db) //Implemente esta funçāo
	postRepo := persistence.NewPostRepository(db)
	log.Println("Repo OK")
	//Serviços de Aplicaçāo
	followService := userApplication.NewFollowService(userRepo)
	postService := postApplication.NewPostService(postRepo, userRepo)
	userService := userApplication.NewUserService(userRepo, postRepo)
	log.Println("Services OK")
	//Handlers HTTP
	userHandler := userApi.NewUserHandlers(followService, userService, postService)
	postHandlers := postApi.NewPostHandlers(postService)
	log.Println("Handler OK")
	// 3 - Configuraçāo do Router
	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)
	setupRoutes(r, userHandler, postHandlers)
	log.Println("Rotas OK")
	log.Println("Iniciando API Social Meli")
	// 4 - Inicializaçāo do Servidor
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
