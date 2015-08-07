package auth

import (
	"fmt"
	"time"

	"gopkg.in/redis.v3"
)

// missingCacheError represents a missing key in the cache
type missingCacheError struct {
	key string
}

// Error formatting for missingCacheError
func (e *missingCacheError) Error() string {
	return fmt.Sprintf("%s", e.key)
}

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type Redis struct {
	client *redis.Client
}

func (r Redis) Get(key string) (string, error) {
	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		return "", &missingCacheError{key}
	} else if err != nil {
		panic(err)
	} else {
		return val, nil
	}
}

func (r Redis) Set(key, value string) error {
	err := r.client.Set(key, value, 10*time.Minute).Err()
	if err != nil {
		panic(err)
	} else {
		return nil
	}
}

func GetRedisCache(hostname, port, password string, database int64) Cache {
	return Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     hostname + ":" + port,
			Password: password,
			DB:       database,
		})}
}
