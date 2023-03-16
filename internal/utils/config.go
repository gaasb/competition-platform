package utils

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"log"
)

var config Config

type Config struct {
	db *sql.DB
}

func Init() {
	err := godotenv.Load(".env")
	dieIf(err)

	config = Config{}

	setupDatabase()
	log.Println("Successfully connected to database!")
}
