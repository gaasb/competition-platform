package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
)

const DB_DRIVER = "pgx"

func GetDB() *sql.DB {
	return config.db
}

func CloseDB() {
	err := config.db.Close()
	if err != nil {
		log.Println(err.Error())
	}
}

func setupDatabase() {
	dsn := dsnAddr()

	db, err := sql.Open(DB_DRIVER, dsn)
	dieIf(err)

	//boil.SetDB(db)

	err = db.Ping()
	dieIf(err)

	config.db = db
	//boil.SetDB(conn.db)
}

func dsnAddr() string {
	//"postgres://login:password@host.docker.internal:5432/database"
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	return dsn
}
