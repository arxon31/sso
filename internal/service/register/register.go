package register

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/arxon31/sso/internal/repo/postgres"
	"github.com/arxon31/sso/pkg/logger"
)

const saltLength = 16

var (
	ErrSomethingWentWrong    = errors.New("something went wrong")
	ErrUserNameAlreadyExists = errors.New("user name already exists")
)

type userStorage interface {
	SaveUser(ctx context.Context, username, passwordHash, salt string) (userID int64, err error)
}

type registerService struct {
	storage userStorage
}

func NewService(storage userStorage) *registerService {
	return &registerService{storage: storage}
}

func (s *registerService) Register(ctx context.Context, username, password string) (userID int64, err error) {
	const op = "service.register.Register"

	passwordHash, salt, err := s.hashPassword(password)
	if err != nil {
		logger.Logger.Error(op, slog.String("error", err.Error()))
		return 0, ErrSomethingWentWrong
	}

	userID, err = s.storage.SaveUser(ctx, username, passwordHash, salt)
	if err != nil {
		logger.Logger.Error(op, slog.String("error", err.Error()))
		switch {
		case errors.Is(err, postgres.ErrUserAlreadyExists):
			return 0, ErrUserNameAlreadyExists
		default:
			return 0, ErrSomethingWentWrong
		}
	}

	return userID, nil
}

func (s *registerService) hashPassword(password string) (string, string, error) {
	salt := make([]byte, saltLength)

	_, err := rand.Read(salt)
	if err != nil {
		return "", "", err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password+string(salt)), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return string(passwordHash), string(salt), nil
}
