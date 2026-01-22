// @title SocialMeli API
// @version 1.0
// @description API Rest para SocialMeli (followers, posts, promoçōes)
// @BasePath /
// @schemes http

package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/RLdAB/API-Social-ML/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"

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

	// 1) middlewares primeiro
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Use(middleware.RedirectSlashes)

	// 2) rota do swagger antes ou depois do setupRoutes (tanto faz),
	// desde que ainda não tenha chamado r.Use depois.
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// 3) rotas da API
	setupRoutes(r, userHandler, postHandlers)
	log.Println("Rotas OK")
	log.Println("Iniciando API Social Meli")
	// 4 - Inicializaçāo do Servidor
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
