package storage

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
)

type ConfirmationStorage interface {
	Create(ctx context.Context, c *models.Confirmation) error
	GetConfirmation(ctx context.Context, name string) (*models.Confirmation, error)
	Delete(ctx context.Context, name string) error
}
