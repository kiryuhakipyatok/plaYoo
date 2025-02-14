package models

import "github.com/google/uuid"

type Notice struct{
	Id 				uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();" json:"notice_id"`
	EventId 		uuid.UUID 	`gorm:"not null;type:uuid"`
	Body 			string 		`json:"body" gorm:"not null"`
}