package controllers

import (
	"avantura/backend/internal/db"
	"avantura/backend/internal/models"
	"math"
	"github.com/gofiber/fiber/v2"
	e "avantura/backend/pkg/error-patterns"
)

func EditRating(c *fiber.Ctx) error{
	userId:=c.Params("id")
	var request struct{
		Stars int `json:"stars"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	user:=models.User{}
	if err:=db.Database.First(&user,"id=?",userId).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	user.NumberOfRatings++
	user.TotalRating += request.Stars
	averageRating := float64(user.TotalRating) / float64(user.NumberOfRatings)
	user.Rating = math.Round(averageRating*2)/2
	if err:=db.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Failed to save changes",
		})
	}
	return c.JSON(fiber.Map{
		"message":"Rating changed",
	})
}
