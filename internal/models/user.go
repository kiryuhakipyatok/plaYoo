package models

import(
	"github.com/lib/pq"
	//"github.com/google/uuid"
) 

type User struct {
	Id 					string 				`json:"id"`
	Login 				string 				`json:"login" gorm:"not null;unique"`
	Tg 					string 				`json:"tg" gorm:"not null"`
	ChatId 				string 				`json:"chat_id" gorm:"uniqe"`
	Followers 			pq.StringArray 		`gorm:"type:text[]" json:"followers"`
	Followings			pq.StringArray 		`gorm:"type:text[]" json:"followings"`
	Rating 				float64 			`json:"rating"`
	TotalRating 		int 				`json:"total_rating"`
	NumberOfRatings 	int 				`json:"num_of_ratings"`
	Events 				pq.StringArray 		`gorm:"type:text[]" json:"events"`
	Comments 			pq.StringArray 		`gorm:"type:text[]" json:"comments"`
	Games 				pq.StringArray 		`gorm:"type:text[]" json:"games"`	
	Notifications		pq.StringArray		`gorm:"type:text[]" json:"notifications"`
	Password 			[]byte 				`json:"-" gorm:"not null"`
	Avatar 				string
	Discord 			string 	
}




