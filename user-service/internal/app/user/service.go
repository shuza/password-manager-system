package user

import (
	"context"
	"fmt"
	"user-service/internal/app/model"
	svc "user-service/internal/app/service"
	"user-service/internal/app/util"
)

type service struct {
	repo svc.UserRepository
}

func NewService(repo svc.UserRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, user model.User) error {
	if user.Email == "" || user.Password == "" || user.FullName == "" || user.BusinessName == "" {
		return fmt.Errorf("invalid user request :%w", model.ErrInvalid)
	}
	if p, err := util.HashPassword(user.Password); err == nil {
		user.Password = p
	}
	return s.repo.InsertUser(ctx, user)
}

func (s *service) GetUserByEmailAndPassword(ctx context.Context, email, password string) (model.User, error) {
	if email == "" || password == "" {
		return model.User{}, fmt.Errorf("invalid login request :%w", model.ErrInvalid)
	}
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	if !util.CheckPasswordHash(password, user.Password) {
		return model.User{}, fmt.Errorf("wrong password :%w", model.ErrInvalid)
	}
	return user, nil
}
