package service

import (
	"context"
	"microBook/internal/domain"
	"microBook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = repository.ErrUserDuplicateEmail
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	//密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	//然后存起来
	return svc.repo.Create(ctx, u)
}
func (svc *UserService) Login(ctx context.Context, email string, password string) error {
	err := svc.repo.FindByEmail(ctx, email)
	return nil
}
