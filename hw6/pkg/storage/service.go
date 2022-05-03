package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(host, port string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // TODO: Add from config
		DB:       0,  // use default DB
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}
	c := &RedisClient{
		Client: client,
	}
	return c, nil
}

func (c *RedisClient) Close() error {
	return c.Client.Close()
}

func (c *RedisClient) GetUser(ctx context.Context, id int) (*models.User, error) {
	data, err := c.Get(ctx, strconv.Itoa(id)).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	u := models.User{}
	err = json.Unmarshal(data, &u)
	if err != nil {
		return nil, fmt.Errorf("can't decode data: %s", err)
	}
	return &u, nil
}
