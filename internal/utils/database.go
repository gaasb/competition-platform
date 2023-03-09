package internal

//go:generate sqlboiler --wipe psql
import (
	"context"
	"database/sql"
	"fmt"
	dbmodels "github.com/gaasb/competition-platform/internal/boiler-models"
	"github.com/gaasb/competition-platform/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
)

const DB_DRIVER = "pgx"

type Postgres struct {
	//conn *sql.Conn
	db *sql.DB
}

var conn Postgres

func InitDB() {

	//"postgres://blank:qwerty@host.docker.internal:5432/main"
	dSN := dbConnAddr()
	db, err := sql.Open(DB_DRIVER, dSN)
	utils.HandleError(err)
	//defer db.Close()
	err = db.Ping()
	utils.HandleError(err)
	fmt.Println("Successfully connected!")
	conn = Postgres{db}
	fmt.Println(dbmodels.Tests().One(context.Background(), db))
	//boil.SetDB(conn.db)
}

func dbConnAddr() string {
	dSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	return dSN
}

func GetRef() *sql.DB {
	return conn.db
}
func GetVal(ctx context.Context) dbmodels.TestSlice {
	tes, _ := dbmodels.Tests().All(ctx, conn.db)
	return tes
}

//CREATE USER docker;
//CREATE DATABASE app;
//GRANT ALL PRIVILEGES ON DATABASE app TO docker;
