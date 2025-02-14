package models

import(
	// "github.com/google/uuid"
	"github.com/lib/pq"
	"time"
) 

type Event struct{
	Id 			  string 			`json:"event_id"`
	AuthorId 	  string 			`gorm:"not null" json:"author_id"`
	Body 		  string 			`json:"body" gorm:"not null"`
	Game 		  string 			`json:"game" gorm:"not null"`
	Members  	  pq.StringArray 	`gorm:"type:text[]" json:"members"`
	Comments 	  pq.StringArray 	`gorm:"type:text[]" json:"comments"`
	Max 		  int 				`json:"max" gorm:"not null"`
	Time 		  time.Time 		`json:"minute" gorm:"not null"`
	NotifiedPre   bool 
}