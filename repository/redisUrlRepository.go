package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

const URL_NOT_FOUND_ERR_MSG_TEMPLATE = "Short URL %s not found"
type ShortURLNotFoundError struct {
	key string
}

func (e ShortURLNotFoundError) Error() string {
	return fmt.Sprintf(URL_NOT_FOUND_ERR_MSG_TEMPLATE, e.key)
}

type redisUrlRepository struct {
    client *redis.Client
}

func NewRedisUrlRepository(db *redis.Client) UrlRepository {
    return &redisUrlRepository{db}
}

func (repo *redisUrlRepository) Store(key string, url string) error {

    err := repo.client.Set(ctx, key, url, 0).Err()
    if err != nil {
        return err
    }

    return nil
}

func (repo *redisUrlRepository) Get(key string) (string, error) {

    url, err := repo.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return "", ShortURLNotFoundError{key}
    } else if err != nil {
        return "", err;
    }

    return url, nil
}