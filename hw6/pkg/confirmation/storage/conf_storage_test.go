package storage_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/confirmation/storage"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
	//// "github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/user/storage"
)

var (
	testClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", "localhost", "6379"),
		Password: "",
		DB:       0, // use default DB
	})
	targetConf = models.Confirmation{
		UserID: 1,
		Code:   "testcode",
	}
	ttl = time.Duration(2 * time.Minute)
)

func TestCreate(t *testing.T) {
	client, err := storage.NewConfirmationStorage(
		"localhost",
		"6379",
		time.Duration(2*time.Minute),
	)
	assert.NoError(t, err)
	err = client.Create(context.Background(), &targetConf)
	assert.NoError(t, err)
	real, err := testClient.Get(context.Background(), "1").Bytes()
	assert.NoError(t, err)
	uReal := models.Confirmation{}
	_ = json.Unmarshal(real, &uReal)
	assert.Equal(t, targetConf, uReal)
	testClient.Del(context.Background(), "1")
}

func TestGet(t *testing.T) {
	testClient.Set(context.Background(), "1", &targetConf, ttl)
	client, err := storage.NewConfirmationStorage(
		"localhost",
		"6379",
		time.Duration(2*time.Minute),
	)
	assert.NoError(t, err)
	res, err := client.GetConfirmation(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, targetConf, *res)
}
