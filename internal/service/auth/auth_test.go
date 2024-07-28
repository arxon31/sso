package auth

import (
	"context"
	"errors"
	"github.com/arxon31/sso/internal/repo/postgres"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
)

func TestAuth_NewService(t *testing.T) {
	provider := &passwordProviderMock{}

	auth := NewService(provider)
	require.NotNil(t, auth)
	require.IsType(t, &authService{}, auth)
}

func TestAuth_Authorize(t *testing.T) {
	err := os.Setenv("SECRET_KEY", "secret")
	require.Nil(t, err)

	t.Run("happy_path", func(t *testing.T) {
		password := "password"
		salt := "salt"
		provider := &passwordProviderMock{
			UserPasswordFunc: func(ctx context.Context, username string) (string, string, error) {

				hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
				require.Nil(t, err)

				return string(hash), salt, nil
			},
		}

		auth := NewService(provider)
		token, err := auth.Authorize(context.Background(), "username", password)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("user_not_exists", func(t *testing.T) {
		provider := &passwordProviderMock{
			UserPasswordFunc: func(ctx context.Context, username string) (string, string, error) {
				return "", "", postgres.ErrUserNotExists
			}}

		auth := NewService(provider)
		token, err := auth.Authorize(context.Background(), "username", "password")
		require.ErrorIs(t, err, ErrWrongLoginOrPassword)
		require.Empty(t, token)
	})

	t.Run("db_error", func(t *testing.T) {
		provider := &passwordProviderMock{
			UserPasswordFunc: func(ctx context.Context, username string) (string, string, error) {
				return "", "", errors.New("some db error")
			}}

		auth := NewService(provider)
		token, err := auth.Authorize(context.Background(), "username", "password")
		require.ErrorIs(t, err, ErrSomethingWentWrong)
		require.Empty(t, token)
	})
}
