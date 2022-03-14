package connector

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type Connector interface {
	AddToEnv(ctx context.Context, user models.User, env models.Env) error
	GetByEnv(ctx context.Context, env models.Env) ([]models.User, error)
	GetByUser(ctx context.Context, user models.User) ([]models.Env, error)
	DeleteUserFromEnv(ctx context.Context, user models.User, env models.Env) error
}

type ConnectorStorage interface {
	AddToEnv(ctx context.Context, user models.User, env models.Env) error
	GetByEnv(ctx context.Context, env models.Env) ([]models.User, error)
	GetByUser(ctx context.Context, user models.User) ([]models.Env, error)
	DeleteUserFromEnv(ctx context.Context, user models.User, env models.Env) error
}
