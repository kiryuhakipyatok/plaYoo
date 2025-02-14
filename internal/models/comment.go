package models

import (
	"time"

	// "github.com/google/uuid"
) 

type Comment struct{
	Id 				string 		`json:"comment_id"`
	AuthorId 		string 		`gorm:"not null" json:"author_id"`
	AuthorName 		string		`gorm:"not null;type:text" json:"author_name"`
	AuthorAvatar 	string		
	Body 			string 		`json:"body" gorm:"not null"`
	Receiver 		string 		`json:"receiver_id"`
	Time 			time.Time	`json:"time" gorm:"not null"`
}