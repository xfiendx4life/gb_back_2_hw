package usecase

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/user"
	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type UserCase struct {
	storage user.UserStorage
}

func New(storage user.UserStorage) user.User {
	return &UserCase{
		storage: storage,
	}
}

func (u *UserCase) Create(ctx context.Context, name string) (user models.User, err error) {
	return u.storage.Create(ctx)
}
func (u *UserCase) GetByName(ctx context.Context, name string) (*models.User, error) {
	return u.storage.GetByName(ctx, name)
}
