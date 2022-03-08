package usecase

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/connector"
	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/env"
	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/user"
	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type Conn struct {
	user      user.User
	env       env.Env
	connStore connector.ConnectorStorage
}

func New(u user.User, en env.Env) connector.Connector {
	return &Conn{
		user: u,
		env:  en,
	}
}

func (c *Conn) AddToEnv(ctx context.Context, user models.User, env models.Env) error {
	return c.connStore.AddToEnv(ctx, user, env)
}
func (c *Conn) AddEnvToUser(ctx context.Context, env models.Env, user models.User) error {
	return c.connStore.AddEnvToUser(ctx, env, user)
}
func (c *Conn) GetByEnv(ctx context.Context, env models.Env) ([]models.User, error) {
	return c.connStore.GetByEnv(ctx, env)
}
func (c *Conn) GetByUser(ctx context.Context, user models.User) ([]models.Env, error) {
	return c.connStore.GetByUser(ctx, user)
}
