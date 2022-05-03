package storage

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
)

type UserStorage interface {
	Create(ctx context.Context, user models.User)
	GetUser(ctx context.Context, id int) (models.User, error)
}
