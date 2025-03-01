package controllers

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
	e "avantura/backend/storage/error-patterns"
)

func GetNotifications(c *fiber.Ctx) error{
	id:=c.Params("id")
	user:=models.User{}
	if err:=postgres.Database.First(&user,"id=?",id).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	return c.JSON(user.Notifications)
}

func DeleteNotification(c *fiber.Ctx) error{
	id:=c.Params("id")
	var request struct{
		Id 		string		`json:"notice_id"` 
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	// userId,_:=uuid.Parse(id)
	user:=models.User{}
	if err:=postgres.Database.First(&user,"id=?",id).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	// noticeId,_:=uuid.Parse(request.Id)
	updateNotifications:=make([]string,0,len(user.Notifications))
	for _,n:=range user.Notifications{
		if n!=request.Id{
			updateNotifications = append(updateNotifications, n)
		}
	}
	user.Notifications = updateNotifications
	if err:=postgres.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Error update user",
		})
	}
	notice:=models.Notice{}
	if err:=postgres.Database.Delete(&notice,"id=?",request.Id).Error;err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"Error deleting notice",
		})
	}
	return c.JSON(fiber.Map{
		"message":"Notification deleted",
	})
	
}


func DeleteAllNotifications(c *fiber.Ctx) error{
	id := c.Params("id")
	user:=models.User{}
	if err:=postgres.Database.First(&user,"id=?",id).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	user.Notifications = nil
	postgres.Database.Save(&user)
	return c.JSON(fiber.Map{
		"message":"All notifications are deleted",
	})
}