package models

import (
	"time"
	"github.com/google/uuid"
) 

type Comment struct{
	Id 				uuid.UUID   	`gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"comment_id"`
	AuthorId 		uuid.UUID 		`gorm:"type:uuid;not null" json:"author_id"`
	AuthorName 		string		`gorm:"not null;type:text" json:"author_name"`
	AuthorAvatar 	string		
	Body 			string 		`json:"body" gorm:"not null"`
	Receiver 		uuid.UUID 		`gorm:"type:uuid" json:"receiver_id"`
	Time 			time.Time	`json:"time" gorm:"not null"`
}