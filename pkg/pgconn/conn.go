package pgconn

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/stdlib"
)

const (
	envPostgresHost = "PG_HOST"
	envPostgresPort = "PG_PORT"
	envPostgresUser = "PG_USER"
	envPostgresPass = "PG_PASS"
	envPostgresDB   = "PG_DB"
)

func New() (*sql.DB, error) {
	postgresDsn := fmt.Sprintf("pgconn://%s:%s@%s:%s/%s?sslmode=disable",
		fromEnv(envPostgresUser),
		fromEnv(envPostgresPass),
		fromEnv(envPostgresHost),
		fromEnv(envPostgresPort),
		fromEnv(envPostgresDB),
	)

	return sql.Open("pgx", postgresDsn)
}

func fromEnv(env string) string {
	val, exists := os.LookupEnv(env)
	if !exists {
		panic(fmt.Sprintf("env %s is not provided", env))
	}

	return val
}
