package repository

import (
	"context"
	"testing"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedisContainer(t *testing.T) (string, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	endpoint, err := redisC.Endpoint(ctx, "")
	assert.NoError(t, err)

	teardown := func() {
		assert.NoError(t, redisC.Terminate(ctx))
	}

	return endpoint, teardown
}

func TestRedisRepository(t *testing.T) {
	endpoint, teardown := setupRedisContainer(t)
	defer teardown()

	repo := &RedisRepository{
		client: redis.NewClient(&redis.Options{
			Addr:     endpoint,
			Password: "",
		}),
	}

	ctx := context.Background()

	t.Run("Test Set and Get", func(t *testing.T) {
		key := "test-key"
		value := "test-value"
		expiration := time.Minute

		err := repo.Set(ctx, key, value, expiration)
		assert.NoError(t, err)

		result, err := repo.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Test Delete", func(t *testing.T) {
		key := "test-key-delete"
		value := "test-value"

		err := repo.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		err = repo.Delete(ctx, key)
		assert.NoError(t, err)

		exists, err := repo.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Test Exists", func(t *testing.T) {
		key := "test-key-exists"
		value := "test-value"

		err := repo.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		exists, err := repo.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)

		err = repo.Delete(ctx, key)
		assert.NoError(t, err)

		exists, err = repo.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Test Close", func(t *testing.T) {
		err := repo.Close()
		assert.NoError(t, err)
	})
}
