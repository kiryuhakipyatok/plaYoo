package postgres

import (
	"avantura/backend/internal/models"
	
	"log"
	
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var Database *gorm.DB

func ConnectToPostgres() error{
	if err := godotenv.Load();err != nil {
        log.Fatal("Error loading .env file"+ err.Error())
    }
	
	var (
		database   = os.Getenv("DB_DATABASE_RLW")
		password   = os.Getenv("DB_PASSWORD_RLW")
		username   = os.Getenv("DB_USERNAME_RLW")
		port       = os.Getenv("DB_PORT_RLW")
		host       = os.Getenv("DB_HOST_RLW")
	)
	if database == "" || username == "" || password == "" || host == "" || port == "" {
        log.Fatal("One or more environment variables are not set")
    }

	dsn:="host=junction.proxy.rlwy.net user=postgres password=ZFvBKMDypzRLBbFqMHMOiRCmFvMPPLCv dbname=railway port=43543 sslmode=disable TimeZone=Europe/Minsk"
	connection,err:= gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		return err
	}
	
	Database = connection
	connection.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Comment{},
		&models.News{},
		&models.Game{},
		&models.Notice{})
	log.Printf("Connect to postgress successfully")
	return nil
}	
