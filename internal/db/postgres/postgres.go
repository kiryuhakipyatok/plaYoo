package postgres

import (
	"avantura/backend/internal/models"
	
	"log"
	

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect(){
	// var (
	// 	database   = os.Getenv("DB_DATABASE_RLW")
	// 	password   = os.Getenv("DB_PASSWORD_RLW")
	// 	username   = os.Getenv("DB_USERNAME_RLW")
	// 	port       = os.Getenv("DB_PORT_RLW")
	// 	host       = os.Getenv("DB_HOST_RLW")
	// )
	dsn:="host=junction.proxy.rlwy.net user=postgres password=ZFvBKMDypzRLBbFqMHMOiRCmFvMPPLCv dbname=railway port=43543 sslmode=disable TimeZone=Europe/Minsk"
	connection,err:= gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		log.Printf("Error connect to database: %v", err)
	}
	Database = connection
	connection.AutoMigrate(&models.User{})
	connection.AutoMigrate(&models.Event{})
	connection.AutoMigrate(&models.Comment{})
	connection.AutoMigrate(&models.News{})
	connection.AutoMigrate(&models.Game{})
	connection.AutoMigrate(&models.Notice{})
}	
