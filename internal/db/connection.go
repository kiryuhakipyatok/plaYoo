package db

import (
	"avantura/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var Database *gorm.DB

func Connect(){
	dsn:="host=dpg-cunourl6l47c73dhnbug-a user=kipyatok password=oaOEuY5zuBdK3jhmaQmHXRFtiULB29Y0 dbname=avanturadb_qgbr port=5432 sslmode=disable TimeZone=Europe/Minsk"
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