package repository

import (
	"context"
	"microBook/internal/domain"
	"microBook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}
func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}
