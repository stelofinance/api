package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stelofinance/api/db"
	"github.com/stelofinance/api/tools"
)

var DB *pgxpool.Pool
var Q *db.Queries

func ConnectDb() error {
	dbpool, err := pgxpool.New(context.Background(), tools.EnvVars.DbConnection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	DB = dbpool
	Q = db.New(dbpool)
	return err
}
