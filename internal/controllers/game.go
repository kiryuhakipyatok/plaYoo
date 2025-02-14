package controllers

import (
	"avantura/backend/internal/db"
	"avantura/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
	e "avantura/backend/pkg/error-patterns"
)

func AddGame(c *fiber.Ctx) error{
	var request struct{
		UserId   string `json:"user_id"`
		GameName string `json:"game"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	user:=models.User{}
	//userId,_:=uuid.Parse(request.UserId)
	if err := db.Database.First(&user,"id=?",request.UserId).Error; err != nil {
        return e.NotFound("user",err,c)
    }
	for _,game:=range user.Games{
		if game == request.GameName{
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message":"The game has already been added",
			})
		}
	}
	game:=models.Game{}
	if err := db.Database.First(&game,"name=?",request.GameName).Error; err != nil {
        c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Game not found",
		})
    }
	user.Games=append(user.Games, request.GameName)
	if err:=db.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Failed to add game",
		})
	}
	game.NumberOfPlayers++
	if err:=db.Database.Save(&game).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Failed to update game",
		})
	}
	

	return c.JSON(fiber.Map{
        "message": "Game added successfully",
    })
}

func DeleteGame(c *fiber.Ctx) error{
	userId:=c.Params("id")
	var request struct{
		Game	string		`json:"game"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	// userIdUUID,_:=uuid.Parse(userId)
	user:=models.User{}
	if err:=db.Database.First(&user,"id=?",userId).Error;err!=nil{
		return e.NotFound("User",err,c)
	}

	updateGames:=make([]string,0,len(user.Games))
	for _, g := range user.Games {
		if g != request.Game {
			updateGames = append(updateGames, g)
		}
	}
	user.Games = updateGames
	if err:=db.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"Failed to update user games",
		})
	}
	
	return c.JSON(fiber.Map{
		"message":"Game deleted successfully",
	})
}

func GetAllGames(c *fiber.Ctx) error{
	games:=[]models.Game{}
	if err:=db.Database.Find(&games).Error;err!=nil{
		return e.ErrorFetching("games",c,err)
	}
	return c.JSON(games)
}


func GetConcreteGame(c *fiber.Ctx) error{
	searchGame:=c.Params("name")
	game:=models.Game{}
	if err := db.Database.First(&game,"name=?",searchGame).Error; err != nil {
        return e.NotFound("Game",err,c)
    }	
	return c.JSON(game)
}