package env

import (
	"context"

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
