package storage

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/env"
	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type someDb struct {
	storage map[string]models.Env
}

// * this is data mapper
type EnvMapper struct {
	db *someDb
}

func New() env.EnvStorage {
	return &EnvMapper{
		db: &someDb{
			storage: make(map[string]models.Env),
		},
	}
}

func (env *EnvMapper) Create(ctx context.Context, name string) (e models.Env, err error) {

	return models.Env{}, nil
}

func (env *EnvMapper) Get(ctx context.Context, name string) (models.Env, error) {
	return models.Env{}, nil
}
