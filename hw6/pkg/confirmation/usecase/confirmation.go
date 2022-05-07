package usecase

import "context"

type Confirmation interface {
	Create(ctx context.Context, name string)
}
