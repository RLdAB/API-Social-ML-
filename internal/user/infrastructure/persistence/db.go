package persistence

import (
	"log"

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
	return db
}
