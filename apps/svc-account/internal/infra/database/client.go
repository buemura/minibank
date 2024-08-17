package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/buemura/minibank/svc-account/config"
	_ "github.com/lib/pq"
)

// var (
// 	Conn *pgxpool.Pool
// )

// func Connect() {
// 	dbConfig, err := pgxpool.ParseConfig(config.DATABASE_URL)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Failed to create pool config: %v\n", err)
// 		os.Exit(1)
// 	}

// 	dbConfig.MaxConns = 100

// 	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	Conn = pool
// }

var (
	Conn *sql.DB
)

func Connect() {
	connectionString :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DATABASE_HOST, config.DATABASE_PORT, config.DATABASE_USER, config.DATABASE_PASS, config.DATABASE_DBNAME)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	Conn = db
}
