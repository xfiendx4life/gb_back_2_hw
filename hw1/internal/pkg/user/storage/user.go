package storage

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/user"
	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type someDb struct {
	storage map[string]models.User
}

type UserStorageMapper struct {
	DB *someDb
}

func New() user.UserStorage {
	return &UserStorageMapper{
		DB: &someDb{
			storage: make(map[string]models.User),
		},
	}
}

func (usm *UserStorageMapper) Create(ctx context.Context) (user models.User, err error) {
	return models.User{}, nil
}
func (usm *UserStorageMapper) GetByName(ctx context.Context, name string) (*models.User, error) {
	return &models.User{}, nil
}
