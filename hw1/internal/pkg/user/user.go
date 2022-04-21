package user

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type UserStorage interface {
	Create(ctx context.Context) (user models.User, err error)
	GetByName(ctx context.Context, name string) (*models.User, error)
}

type User interface {
	Create(ctx context.Context, name string) (user models.User, err error) // * actually creates user and call storage layer
	GetByName(ctx context.Context, name string) (*models.User, error)
}

type UserDeliver interface {
	Create(ectx echo.Context) error
	GetByName(ectx echo.Context) error
}
