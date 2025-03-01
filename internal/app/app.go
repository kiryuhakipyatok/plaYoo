package app

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/db/redis"
	"avantura/backend/internal/notify"
	"avantura/backend/internal/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)
func Run() {
	var wg sync.WaitGroup
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
	stop:=make(chan struct{})
	wg.Add(1)
	go func(){
		notify.CreateBot(stop)
		defer wg.Done()
	}()
	wg.Add(1)
	go func ()  {
		notify.ScheduleNotify(stop)
		defer wg.Done()
	}()

	app:=server.RunServer()
	
	quit:=make(chan os.Signal,1)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	time.Sleep(3*time.Second)
	if err:=app.Shutdown();err!=nil{
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	close(stop)
	wg.Wait()
	log.Println("Server stopped")
}