package controllers

import (
	"avantura/backend/internal/db"
	"avantura/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	e "avantura/backend/pkg/error-patterns"
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
	userIdUUID,_:=uuid.Parse(request.UserId)
	if err:=db.Database.First(&user,"id=?",userIdUUID).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	user.Discord = request.Discord
	if err:=db.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Error save changes",
		})
	}
	return c.JSON(fiber.Map{
		"message":"Discord saved",
	})
}