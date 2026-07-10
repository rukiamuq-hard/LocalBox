package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const passwordRedis = "mypass"

type RedisDB struct {
	rdb *redis.Client
}

func New() *RedisDB {
	return &RedisDB{}
}

func (myRDB *RedisDB) CreateRedis() {
	myRDB.rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6380",
		Password: "",
		DB:       0,
	})
}

func (myRDB *RedisDB) SetKeyValue(key string, value int) error {
	err := myRDB.rdb.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (myRDB *RedisDB) GetValue(key string) (string, error) {
	val, err := myRDB.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (myRDB *RedisDB) Close() {
	myRDB.rdb.Close()
}
