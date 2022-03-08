package usecase

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/env"
	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type Env struct {
	storage env.EnvStorage
}

func New(storage env.EnvStorage) env.Env {
	return &Env{
		storage: storage,
	}
}

func (e *Env) Create(ctx context.Context, name string) (env models.Env, err error) {
	return e.storage.Create(ctx, name)
}
func (e *Env) Get(ctx context.Context, name string) (env models.Env, err error) {
	return e.storage.Get(ctx, name)
}
