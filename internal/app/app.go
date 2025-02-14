package app

import (
	"avantura/backend/internal/db"
	// "avantura/backend/internal/notify"
	"avantura/backend/internal/server"
	// "sync"
)
func Run() {
    db.Connect()
	// go notify.CreateBot()
	server.RunServer()
    // var wg sync.WaitGroup
    // wg.Add(1)
    // wg.Wait()
}