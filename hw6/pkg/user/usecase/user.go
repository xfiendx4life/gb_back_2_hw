package usecase

import "context"

type User interface {
	Register(ctx context.Context, name, password string) error
	Confirm(ctx context.Context, name string) error
}
