package app

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/db/redis"
	"avantura/backend/internal/notify"
	"avantura/backend/internal/server"
	"log"
	"os/signal"
	"os"
	"syscall"
	"time"
)
func Run() {
	if err:=postgres.ConnectToPostgres();err!=nil{
		log.Fatalf("Error to connenct to Postgres: %v",err)
	}
	errRedis:=redis.ConnectToRedis()
	if errRedis!=nil{
		log.Printf("Error to connenct to Redis: %v",errRedis)
	}
	closeDB,err:=postgres.Database.DB()
	if err!=nil{
		log.Fatalf("Failed to get DB: %v", err)
	}
	defer func ()  {
		if err:=closeDB.Close();err!=nil{
			log.Printf("Error to close Postgres: %v",err)
		}else{
			log.Printf("Close Postgres success")
		}
		if errRedis==nil{
			if err:=redis.Rdb.Close();err!=nil{
				log.Printf("Error to close Redis: %v",err)
			}else{
				log.Printf("Close Redis success")
			}
		}
	}() 
	//defer redis.Rdb.Close()
	go notify.CreateBot()
	go notify.ScheduleNotify()
	app:=server.RunServer()
	quit:=make(chan os.Signal,1)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	time.Sleep(3*time.Second)
	if err:=app.Shutdown();err!=nil{
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
    // var wg sync.WaitGroup
    // wg.Add(1)
    // wg.Wait()
}