package postgres

import (
	"avantura/backend/internal/models"
	"fmt"

	"log"

	// "github.com/joho/godotenv"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectToPostgres() error{
	// if err := godotenv.Load();err != nil {
    //     log.Fatal("Error loading .env file"+ err.Error())
    // }
	
	var (
		database   = os.Getenv("POSTGRES_DB")
		password   = os.Getenv("POSTGRES_PASSWORD")
		username   = os.Getenv("POSTGRES_USER")
		port       = os.Getenv("PGPORT")
		host       = os.Getenv("PGHOST")
	)
	if database == "" || username == "" || password == "" || host == "" || port == "" {
        log.Fatal("One or more environment variables are not set")
    }
	dsn:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Minsk",host,username,password,database,port)
	//dsn:="host=junction.proxy.rlwy.net user=postgres password=ZFvBKMDypzRLBbFqMHMOiRCmFvMPPLCv dbname=railway port=43543 sslmode=disable TimeZone=Europe/Minsk"
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
