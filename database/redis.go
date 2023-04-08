package database

import (
	"github.com/go-redis/redis"
)

func NewRedisDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "qwerty",
		DB:       0,
	})

	err := rdb.Ping().Err()
	if err != nil {
		panic(err)
	}

	return rdb
}
