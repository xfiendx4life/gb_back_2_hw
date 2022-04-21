package env

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type Env interface {
	Create(ctx context.Context, name string) (env models.Env, err error)
	Get(ctx context.Context, name string) (env models.Env, err error)
}

// * interface for dataMapper
type EnvStorage interface {
	Create(ctx context.Context, name string) (env models.Env, err error)
	Get(ctx context.Context, name string) (env models.Env, err error)
}

type EnvDeliver interface {
	Create(ectx echo.Context) error
	Get(ectx echo.Context) error
}
