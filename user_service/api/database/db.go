package database

import (
	"context"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
)

// DB postgres db client
var DB *pg.DB

// Connect creates a connection to postgres
func Connect() error {

	opt, err := pg.ParseURL(os.Getenv("DB_URI_HEROKU"))
	if err != nil {
		return err
	}

	DB = pg.Connect(opt)

	// Check if database is up and running
	if err := DB.Ping(context.Background()); err != nil {
		return err
	}
	log.Println("Connected to db")
	return nil
}
