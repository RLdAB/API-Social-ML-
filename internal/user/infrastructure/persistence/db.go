package persistence

import (
	"fmt"
	"log"
	"os"
	"time"

	postdomain "github.com/RLdAB/API-Social-ML/internal/post/domain"
	userdomain "github.com/RLdAB/API-Social-ML/internal/user/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// getEnv retorna a variável de ambiente ou um default.
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// NewDB inicializa e retorna a conexão com o banco e executa as migrations.
// Funciona tanto local (host=localhost) quanto via docker-compose (host=postgres).
func NewDB() *gorm.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "socialmeli")
	pass := getEnv("DB_PASSWORD", "socialmeli")
	name := getEnv("DB_NAME", "socialmeli")
	ssl := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host, user, pass, name, port, ssl,
	)

	// Pequeno retry (ajuda muito no Docker, quando o Postgres ainda está subindo)
	var db *gorm.DB
	var err error
	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("tentativa %d/10: falha ao conectar ao banco (%v)", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("falha ao conectar ao banco: %v", err)
	}

	// Migração das tabelas:
	if err := db.AutoMigrate(
		&userdomain.User{},
		&userdomain.Follow{},
		&postdomain.Post{},
	); err != nil {
		log.Fatalf("falha ao executar AutoMigrate: %v", err)
	}

	log.Println("Banco conectado e migrações executadas com sucesso")
	return db
}