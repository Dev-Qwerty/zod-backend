package database

import (
	"context"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
)

// DB postgres db client
var DB *pg.DB

// Connect creates a connection to postgres
func Connect() error {

	var err error

	err = godotenv.Load()
	if err != nil {
		return err
	}

	DB = pg.Connect(&pg.Options{
		Addr:     os.Getenv("DB_ADDR"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	})

	// Check if database is up and running
	if err := DB.Ping(context.Background()); err != nil {
		return err
	}
	log.Println("Connected to db")
	return nil
}
