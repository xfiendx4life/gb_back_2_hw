package usecase_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/confirmation/usecase"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
)

type mockSt struct {
	res models.Confirmation
	err error
}

var target = models.Confirmation{
	UserName: "testuser",
	Code:     "1234",
}

func (mc *mockSt) Create(ctx context.Context, c *models.Confirmation) error {
	return mc.err
}

func (mc *mockSt) GetConfirmation(ctx context.Context, name string) (*models.Confirmation, error) {
	return &mc.res, mc.err
}
func (mc *mockSt) Delete(ctx context.Context, name string) error {
	return mc.err
}

func TestCreate(t *testing.T) {
	mc := mockSt{res: target}
	conf := usecase.New(&mc)
	code, err := conf.Create(context.Background(), target.UserName)
	assert.NoError(t, err)
	c, err := strconv.Atoi(code)
	assert.NoError(t, err)
	assert.True(t, c < 1000 && c > 99)

}

func TestCreateError(t *testing.T) {
	mc := mockSt{err: fmt.Errorf("testerror")}
	conf := usecase.New(&mc)
	code, err := conf.Create(context.Background(), target.UserName)
	assert.Error(t, err)
	assert.Equal(t, "", code)

}

func TestConfirm(t *testing.T) {
	mc := mockSt{res: target}
	conf := usecase.New(&mc)
	f, err := conf.Confirm(context.Background(), target.UserName,
		target.Code)
	assert.NoError(t, err)
	assert.True(t, f)

}

func TestConfirmFalse(t *testing.T) {
	mc := mockSt{res: target}
	conf := usecase.New(&mc)
	f, err := conf.Confirm(context.Background(), target.UserName,
		"1111")
	assert.NoError(t, err)
	assert.False(t, f)
}

func TestConfirmErr(t *testing.T) {
	mc := mockSt{res: target, err: fmt.Errorf("testerror")}
	conf := usecase.New(&mc)
	f, err := conf.Confirm(context.Background(), target.UserName,
		"1111")
	assert.Error(t, err)
	assert.False(t, f)
}


