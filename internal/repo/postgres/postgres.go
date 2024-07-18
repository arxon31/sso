package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrUserNotExists     = errors.New("user not exists")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) (*postgres, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgres{db: db}, nil
}

func (pg *postgres) SaveUser(ctx context.Context, username, passwordHash, salt string) (userID int64, err error) {
	const op = "repo.pgconn.RegisterUser"

	exists, err := pg.isExists(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	if exists {
		return 0, ErrUserAlreadyExists
	}

	query := "INSERT into users (username, password_hash, salt) VALUES ($1, $2, $3);"

	res, err := pg.db.ExecContext(ctx, query, username, passwordHash, salt)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return res.LastInsertId()
}

func (pg *postgres) UserPassword(ctx context.Context, username string) (passwordHash, salt string, err error) {
	const op = "repo.pgconn.UserPassword"

	query := "SELECT (password_hash, salt) from users WHERE username=$1;"

	row := pg.db.QueryRowContext(ctx, query, username)

	err = row.Scan(&passwordHash, &salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrUserNotExists
		}
		return "", "", err
	}

	return passwordHash, salt, nil
}

func (pg *postgres) isExists(ctx context.Context, username string) (bool, error) {
	const op = "repo.pgconn.isExists"

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1);"

	row := pg.db.QueryRowContext(ctx, query, username)

	var exists bool

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return exists, nil
}
