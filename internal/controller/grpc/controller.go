package grpc

import (
	"context"

	"github.com/arxon31/yapr-proto/pkg/sso"
)

type authService interface {
	Authorize(ctx context.Context, username, password string) (token string, err error)
}

type registerService interface {
	Register(ctx context.Context, username, password string) (userID int64, err error)
}

type controller struct {
	auth     authService
	register registerService
	sso.UnimplementedSSOServer
}

func NewController(auth authService, register registerService) *controller {
	return &controller{auth: auth, register: register}
}

func (s *controller) Register(ctx context.Context, request *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	id, err := s.register.Register(ctx, request.GetUsername(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &sso.RegisterResponse{UserId: id}, nil
}

func (s *controller) Login(ctx context.Context, request *sso.LoginRequest) (*sso.LoginResponse, error) {
	token, err := s.auth.Authorize(ctx, request.GetUsername(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &sso.LoginResponse{Token: token}, nil
}

func (s *controller) Logout(ctx context.Context, request *sso.LogoutRequest) (*sso.LogoutResponse, error) {
	return &sso.LogoutResponse{Success: true}, nil
}
