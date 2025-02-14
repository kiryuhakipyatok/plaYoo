package models

import "gorm.io/gorm"

type Game struct{
	Name 				string 		`json:"name" gorm:"not null; primaryKey; unique"`
	NumberOfPlayers 	int 		`json:"num_of_players"`
	NumberOfEvents 		int 		`json:"num_of_events"`
	Rating 				float64 	`json:"rating"`
}

func (g *Game) CalculateRating() {
    if g.NumberOfPlayers == 0 {
        g.Rating = 0 
        return
    }
    g.Rating = float64(g.NumberOfPlayers+g.NumberOfEvents) / 2.0
}

func (g *Game) BeforeUpdate(tx *gorm.DB) (err error) {
    g.CalculateRating()
    return
}