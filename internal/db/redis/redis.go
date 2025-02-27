package redis

import (
	"context"
	"log"
	"os"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx context.Context
)

func ConnectToRedis() error{

	var (
		host = os.Getenv("REDISHOST")
		port = os.Getenv("REDISPORT")
		password = os.Getenv("REDISPASSWORD")
	)

	// client := redis.NewClient(&redis.Options{
    //     Addr:     "localhost:6379",
    //     Password: "",               
    //     DB:       0,              
    // })
	addr:=fmt.Sprintf("%s:%s",host,port)
	client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,               
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