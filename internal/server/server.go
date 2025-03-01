package server

import (
	"avantura/backend/internal/routes"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
)
func RunServer() *fiber.App{
    app:=fiber.New()    
	app.Static("pkg","../../pkg")
    app.Use(cors.New(cors.Config{        
		AllowOrigins:     "*",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",         
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization", 
        ExposeHeaders:    "Content-Length",        
		AllowCredentials: false, 
    }))    
	routes.Setup(app)
	port:=os.Getenv("PORT")
	if port == ""{
		port = "3000"
	}
    log.Fatal(app.Listen("0.0.0.0:"+port))    
	return app
}
