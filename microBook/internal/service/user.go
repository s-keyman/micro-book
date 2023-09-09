package service

import (
	"context"
	"errors"
	"microBook/internal/domain"
	"microBook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("邮箱/密码不对")
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
func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		//这里可以打日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
