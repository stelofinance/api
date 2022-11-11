package database

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/stelofinance/api/db"
	"github.com/stelofinance/api/tools"
)

var DB *pgx.Conn
var Q *db.Queries

func ConnectDb() error {
	conn, err := pgx.Connect(context.Background(), tools.EnvVars.DbConnection)
	DB = conn
	Q = db.New(conn)
	return err
}