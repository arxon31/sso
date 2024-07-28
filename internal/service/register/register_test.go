package register

import (
	"context"
	"errors"
	"github.com/arxon31/sso/internal/repo/postgres"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegister_NewService(t *testing.T) {
	storage := &userStorageMock{}
	s := NewService(storage)
	require.NotNil(t, s)
	require.IsType(t, &registerService{}, s)
}

func TestRegister_Register(t *testing.T) {
	t.Run("happy_path", func(t *testing.T) {
		storage := &userStorageMock{
			SaveUserFunc: func(ctx context.Context, username string, passwordHash, salt []byte) (userID int64, err error) {
				return 1, nil
			},
		}

		s := NewService(storage)
		userID, err := s.Register(context.Background(), "username", "password")
		require.NoError(t, err)
		require.Equal(t, int64(1), userID)
	})

	t.Run("user_already_exists", func(t *testing.T) {
		storage := &userStorageMock{
			SaveUserFunc: func(ctx context.Context, username string, passwordHash, salt []byte) (userID int64, err error) {
				return 0, postgres.ErrUserAlreadyExists
			},
		}

		s := NewService(storage)
		userID, err := s.Register(context.Background(), "username", "password")
		require.ErrorIs(t, err, ErrUserNameAlreadyExists)
		require.Equal(t, int64(0), userID)
	})

	t.Run("error", func(t *testing.T) {
		storage := &userStorageMock{
			SaveUserFunc: func(ctx context.Context, username string, passwordHash, salt []byte) (userID int64, err error) {
				return 0, errors.New("something went wrong")
			},
		}

		s := NewService(storage)
		userID, err := s.Register(context.Background(), "username", "password")
		require.ErrorIs(t, err, ErrSomethingWentWrong)
		require.Equal(t, int64(0), userID)
	})
}
