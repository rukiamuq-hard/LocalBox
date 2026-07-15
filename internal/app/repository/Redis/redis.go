package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const LifeTimeObject = 20

type RedisDB struct {
	rdb *redis.Client
}

func New() *RedisDB {
	return &RedisDB{}
}

func (myRDB *RedisDB) StartRedis() error {
	myRDB.rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err := myRDB.rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

func (myRDB *RedisDB) SetKeyValue(key string, value int) error {
	err := myRDB.rdb.Set(context.Background(), key, value, LifeTimeObject*time.Minute).Err()
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
