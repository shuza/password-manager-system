package service

import (
	"context"
	"user-service/internal/app/model"
)

//go:generate mockgen -source service.go -destination ./mocks/mock_service.go -package mocks

type UserService interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUserByEmailAndPassword(ctx context.Context, email, password string) (model.User, error)
}

type UserRepository interface {
	InsertUser(ctx context.Context, user model.User) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type AuthService interface {
	Decode(tokenStr string) (*model.JwtCustomClaims, error)
	Encode(user model.User) (string, error)
}
