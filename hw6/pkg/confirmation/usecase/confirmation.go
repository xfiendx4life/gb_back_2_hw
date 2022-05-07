package usecase

import "context"

type Confirmation interface {
	Create(ctx context.Context, userName string) (code string, err error)
	Confirm(ctx context.Context, userName, code string) (bool, error)
}
