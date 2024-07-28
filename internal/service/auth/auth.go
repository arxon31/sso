package auth

import (
	"context"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/arxon31/sso/internal/repo/postgres"
	"github.com/arxon31/sso/pkg/jwt"
	"github.com/arxon31/sso/pkg/logger"
)

var (
	ErrWrongLoginOrPassword = errors.New("wrong login or password")
	ErrSomethingWentWrong   = errors.New("something went wrong")
)

//go:generate moq -out passwordProvider_moq_test.go . passwordProvider
type passwordProvider interface {
	UserPassword(ctx context.Context, username string) (passwordHash, salt string, err error)
}

type authService struct {
	pwdProvider passwordProvider
}

func NewService(provider passwordProvider) *authService {
	return &authService{pwdProvider: provider}
}

func (s *authService) Authorize(ctx context.Context, username, password string) (token string, err error) {
	const op = "service.auth.Authorize"
	passHash, salt, err := s.pwdProvider.UserPassword(ctx, username)
	if err != nil {
		logger.Logger.Error(op, slog.String("error", err.Error()))

		switch {
		case errors.Is(err, postgres.ErrUserNotExists):
			return "", ErrWrongLoginOrPassword
		default:
			return "", ErrSomethingWentWrong
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password+salt))
	if err != nil {
		logger.Logger.Error(op, slog.String("error", err.Error()))
		return "", ErrWrongLoginOrPassword
	}

	tokenString, err := jwt.NewToken(username)
	if err != nil {
		logger.Logger.Error(op, slog.String("error", err.Error()))
		return "", ErrSomethingWentWrong
	}

	return tokenString, nil
}
