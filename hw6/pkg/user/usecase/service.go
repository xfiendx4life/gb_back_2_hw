package usecase

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/user/storage"
)

type UserCase struct {
	store storage.UserStorage
}

func NewUserCase(st storage.UserStorage) User {
	// TODO: Send to config
	// host := os.Getenv("HOST")
	// port := os.Getenv("PORT")
	// ttl, err := strconv.Atoi(os.Getenv("TTL"))
	// if err != nil {
	// 	log.Printf("can't parse ttl string %s", err)
	// 	ttl = 1
	// }
	// //st, err := storage.NewUserStorage(host, port, time.Duration(ttl*int(time.Minute)))
	return &UserCase{
		store: st,
	}
}

// Register new user waiting for confirmation
func (u *UserCase) Register(ctx context.Context, name, password string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("register done with context")
	default:
		user := models.User{
			Name:      name,
			Password:  password,
			Confirmed: false,
		}
		err := u.store.Create(ctx, &user)
		if err != nil {
			return fmt.Errorf("can't create user: %s", err)
		}
		return nil
	}
}

func (u *UserCase) Confirm(ctx context.Context, name string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("confirm done with context")
	default:
		us, err := u.store.GetUser(ctx, name)
		if err != nil {
			return fmt.Errorf("can't confirm user: %s", err)
		}
		us.Confirmed = true
		log.Infof("%v", us)
		err = u.store.Create(ctx, us)
		if err != nil {
			return fmt.Errorf("can't set new value to user %s: %s", us.Name, err)
		}
		return nil
	}
}
