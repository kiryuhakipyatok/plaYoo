package controllers

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	e "avantura/backend/storage/error-patterns"
)

func RecordDiscord(c *fiber.Ctx) error {
	var request struct{
		UserId 			string 		`json:"user_id"`
		Discord 		string 		`json:"discord"`
	}
	if err:=c.BodyParser(request.Discord);err!=nil{
		return e.BadRequest(c,err)
	}
	user:=models.User{}
	if err:=postgres.Database.First(&user,"id=?",request.UserId).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	user.Discord = request.Discord
	if err:=postgres.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Error save changes",
		})
	}
	return c.JSON(fiber.Map{
		"message":"Discord saved",
	})
}