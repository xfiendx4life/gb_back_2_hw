package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/user/usecase"
)

type mockStorage struct {
	err error
}

func (ms *mockStorage) Create(ctx context.Context, user *models.User) error {
	return ms.err
}

func (ms *mockStorage) GetUser(ctx context.Context, name string) (*models.User, error) {
	return &models.User{
		Name:     "testname",
		Password: "testpwd",
	}, ms.err
}

func TestRegister(t *testing.T) {
	uc := usecase.NewUserCase(&mockStorage{})
	err := uc.Register(context.Background(), "testname", "testpwd")
	assert.NoError(t, err)
}

func TestRegisterError(t *testing.T) {
	uc := usecase.NewUserCase(&mockStorage{err: fmt.Errorf("test err")})
	err := uc.Register(context.Background(), "testname", "testpwd")
	assert.Error(t, err)
}

func TestRegisterCtx(t *testing.T) {
	uc := usecase.NewUserCase(&mockStorage{})
	ctx, cancelFunc := context.WithCancel(context.Background())
	errChan := make(chan error)
	go func() {
		err := uc.Register(ctx, "testname", "testpwd")
		errChan <- err
	}()
	cancelFunc()
	err := <-errChan
	assert.Error(t, err)
	assert.Equal(t, "register done with context", err.Error())

}

func TestConfirm(t *testing.T) {
	uc := usecase.NewUserCase(&mockStorage{})
	err := uc.Confirm(context.Background(), "testname")
	assert.NoError(t, err)
}

func TestConfirmError(t *testing.T) {
	uc := usecase.NewUserCase(&mockStorage{err: fmt.Errorf("test err")})
	err := uc.Confirm(context.Background(), "testname")
	assert.Error(t, err)
}

func TestConfirmCtx(t *testing.T) {
	uc := usecase.NewUserCase(&mockStorage{})
	ctx, cancelFunc := context.WithCancel(context.Background())
	errChan := make(chan error)
	go func() {
		err := uc.Confirm(ctx, "testname")
		errChan <- err
	}()
	cancelFunc()
	err := <-errChan
	assert.Error(t, err)
	assert.Equal(t, "confirm done with context", err.Error())

}
