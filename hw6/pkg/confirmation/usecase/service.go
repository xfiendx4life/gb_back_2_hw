package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/confirmation/storage"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
)

type confCase struct {
	store storage.ConfirmationStorage
}

func New(st storage.ConfirmationStorage) Confirmation {
	return &confCase{
		store: st,
	}
}

func genCode(a, b int) string {
	rand.Seed(time.Now().Unix())
	return strconv.Itoa(rand.Intn(b-a) + a)
}

func (c *confCase) Create(ctx context.Context, userName string) (code string, err error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("create done with context")
	default:
		conf := models.Confirmation{
			UserName: userName,
			Code:     genCode(100, 999),
		}
		err = c.store.Create(ctx, &conf)
		if err != nil {
			return "", fmt.Errorf("can't create confirmation: %s", err)
		}
		return conf.Code, nil
	}
}

func (c *confCase) Confirm(ctx context.Context, userName, code string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, fmt.Errorf("confirmation done with context")
	default:
		conf, err := c.store.GetConfirmation(ctx, userName)
		if err != nil {
			return false, fmt.Errorf("cant confitm: %s", err)
		}
		if conf.Code != code {
			return false, nil
		}
		err = c.store.Delete(ctx, userName)
		if err != nil {
			return false, fmt.Errorf("can't confirm: %s", err)
		}
		return true, nil
	}
}
