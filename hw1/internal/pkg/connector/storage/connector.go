package storage

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/connector"
	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/env"
	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/user"
	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type Conn struct {
	user user.UserStorage
	env  env.EnvStorage
}

func New(u user.UserStorage, en env.EnvStorage) connector.Connector {
	return &Conn{
		user: u,
		env:  en,
	}
}

func (c *Conn) AddToEnv(ctx context.Context, user models.User, env models.Env) error {
	return nil
}
func (c *Conn) AddEnvToUser(ctx context.Context, env models.Env, user models.User) error {
	return nil
}
func (c *Conn) GetByEnv(ctx context.Context, env models.Env) ([]models.User, error) {
	return []models.User{}, nil
}

func (c *Conn) GetByUser(ctx context.Context, user models.User) ([]models.Env, error) {
	return []models.Env{}, nil
}
