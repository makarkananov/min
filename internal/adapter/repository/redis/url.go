package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type URLRepository struct {
	client *redis.Client
}

func NewURLRepository(client *redis.Client) *URLRepository {
	return &URLRepository{
		client: client,
	}
}

func (r *URLRepository) GetOriginal(ctx context.Context, short string) (string, error) {
	original, err := r.client.Get(ctx, short).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", err
	}

	return original, nil
}

func (r *URLRepository) Add(ctx context.Context, short, original string) error {
	_, err := r.client.Set(ctx, short, original, 0).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *URLRepository) Remove(ctx context.Context, short string) error {
	_, err := r.client.Del(ctx, short).Result()
	if err != nil {
		return err
	}

	return nil
}
