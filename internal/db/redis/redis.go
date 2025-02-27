package redis

import (
	"github.com/redis/go-redis/v9"
	"log"
	"context"
)

var (
	Rdb *redis.Client
	Ctx context.Context
)

func ConnectToRedis() error{
	client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",               
        DB:       0,              
    })
	// defer client.Close()
	Rdb = client
	Ctx = context.Background()
	if _,err:= client.Ping(Ctx).Result();err!=nil{
		return err
	}else{
		log.Printf("Connect to redis successfully")
	}
	return nil
}