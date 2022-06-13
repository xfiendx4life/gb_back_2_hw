package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/models"
)

type Storage interface {
	Create(ctx context.Context, list models.List) error
	Read(ctx context.Context, id uuid.UUID) (list *models.List, err error)
	Update(ctx context.Context, id uuid.UUID, items []*models.Item) error
	Delete(ctx context.Context, id uuid.UUID) error
}
