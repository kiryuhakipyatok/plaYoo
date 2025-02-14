package server

import (
	"avantura/backend/internal/routes"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)
func RunServer() *fiber.App{
    app:=fiber.New()    
	app.Static("pkg","../../pkg")
    app.Use(cors.New(cors.Config{        
		AllowOrigins:     "http://localhost:3000,http://localhost:9110",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",         
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization", 
        ExposeHeaders:    "Content-Length",        
		AllowCredentials: true, 
    }))    
	routes.Setup(app)
	port:=os.GETENV("PORT")
	if port == ""{
		port = "3000"
	}
    log.Fatal(app.Listen(":"+port))    
	return app
}
