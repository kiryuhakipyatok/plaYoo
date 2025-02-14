package models

// import "github.com/google/uuid"

type Notice struct{
	Id 				string		`json:"notice_id"`
	EventId 		string 		`gorm:"not null"`
	Body 			string 		`json:"body" gorm:"not null"`
}