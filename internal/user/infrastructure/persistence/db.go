package persistence

import (
	"log"

	"github.com/RLdAB/API-Social-ML/internal/user/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB inicializa e retorna a conexāo com o banco de dados
func NewDB() *gorm.DB {
	dsn := "host=postgres user=socialmeli password=socialmeli dbname=socialmeli port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("falha ao conectar ao banco: %v", err)
	}

	// Migraçāo AQUI:
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("falha ao fazer AutoMigrate User: %v", err)
	}
	return db
}
