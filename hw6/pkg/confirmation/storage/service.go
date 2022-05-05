package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
)

type confStore struct {
	*redis.Client
	ttl time.Duration
}

func NewConfirmationStorage(host, port string, ttl time.Duration) (ConfirmationStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // TODO: Add from config
		DB:       0,  // use default DB
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}
	c := &confStore{
		Client: client,
	}
	return c, nil
}

func (c *confStore) GetConfirmation(ctx context.Context, userId int) (*models.Confirmation, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetConfirmation done with context")
	default:
		data, err := c.Get(ctx, strconv.Itoa(userId)).Bytes()
		if err == redis.Nil {
			// we got empty result, it's not an error
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		con := models.Confirmation{}
		err = json.Unmarshal(data, &con)
		if err != nil {
			return nil, fmt.Errorf("can't decode data: %s", err)
		}
		return &con, nil
	}

}

func (c *confStore) Create(ctx context.Context, con *models.Confirmation) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("confirmation.Create done with context")
	default:
		err := c.Set(ctx, strconv.Itoa(con.UserID), con, c.ttl).Err()
		if err != nil {
			return fmt.Errorf("can't add data to redis: %s", err)
		}
		return nil
	}
}

func (c *confStore) Delete(ctx context.Context, userID int) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("confirmation delete done with context")
	default:
		err := c.Del(ctx, strconv.Itoa(userID)).Err()
		if err != nil {
			return fmt.Errorf("can't delete confirmation: %s", err)
		}
		return nil
	}
}
