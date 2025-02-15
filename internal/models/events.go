package models

import(
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
) 

type Event struct{
	Id 			  uuid.UUID 		`gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"event_id"`
	AuthorId 	  uuid.UUID			`gorm:"type:uuid;not null" json:"author_id"`
	Body 		  string 			`json:"body" gorm:"not null"`
	Game 		  string 			`json:"game" gorm:"not null"`
	Members  	  pq.StringArray 	`gorm:"type:uuid[]" json:"members"`
	Comments 	  pq.StringArray 	`gorm:"type:uuid[]" json:"comments"`
	Max 		  int 				`json:"max" gorm:"not null"`
	Time 		  time.Time 		`json:"minute" gorm:"not null"`
	NotifiedPre   bool 
}