package redis

import (
	"github.com/redis/go-redis/v9"
	"context"
)

var (
	Rdb *redis.Client
	Ctx context.Context
)

func ConnectToRedis(){
	client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6666",           
        DB:       0,              
    })
	Rdb = client
	Ctx = context.Background()
}