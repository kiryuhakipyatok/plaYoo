package models

import(
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type News struct{
	Id 			uuid.UUID 				`gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"news_id"`
	AuthorId 	uuid.UUID				`gorm:"type:uuid;not null" json:"author_id"`
	AuthorAvatars string 				`gorm:"not null" json:"author_avatar"`
	AuthorName 	string 					`gorm:"not null" json:"author_name"`
	Title 		string 					`gorm:"not null" json:"title"`
	Body 		string 					`json:"body" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Time 		time.Time 				`json:"time" gorm:"not null"`
	Link  		string 					`gorm:"type:text" json:"link"`
	Comments 	pq.StringArray 			`gorm:"type:uuid[]" json:"comments"`
	Picture 	string					`json:"picture"`
}