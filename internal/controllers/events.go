package controllers

import (
	"avantura/backend/internal/db"
	"avantura/backend/internal/models"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	e "avantura/backend/pkg/error-patterns"
)

func AddEvent(c *fiber.Ctx) error{
	iventdata:=map[string]string{}
	if err:=c.BodyParser(&iventdata);err!=nil{
		return e.BadRequest(c,err)
	}
	author_id:=c.Params("id")
	authorid,err:=uuid.Parse(author_id)
	if err != nil {
		return e.BadUUID(c,err)
	}
	user:=models.User{}
	if err:=db.Database.First(&user,"id=?",author_id).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	max,_:=strconv.Atoi(iventdata["max"])
	minute,_:=strconv.Atoi(iventdata["minute"])
	members:=[]string{author_id}
	event:=models.Event{
		Id: uuid.New(),
		AuthorId: authorid,
		Body: iventdata["body"],
		Game:iventdata["game"],
		Max: max,
		Members: members,
		Time: time.Now().Add(time.Duration(minute)*time.Minute),
	}
	if minute < 10 {
		event.NotifiedPre = true
	}
	if err:=db.Database.Create(&event).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error creating ivent",
		})
	}
	game:=models.Game{}
	if err := db.Database.First(&game,"name=?",iventdata["game"]).Error; err != nil {
        return e.NotFound("Game",err,c)
    }
	user.Events=append(user.Events, event.Id.String())
	game.NumberOfEvents++
	db.Database.Save(&user)
	db.Database.Save(&game)
	return c.JSON(event)
}

func GetConcreteEvent(c *fiber.Ctx) error{
	userId:=c.Params("id")
	user:=models.User{}
	if err:=db.Database.First(&user,"id=?",userId).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	return c.JSON(user.Events)
}

func GetEvents(c *fiber.Ctx) error{
	var request struct{
		Amount string `json:"amount"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	a,_:=strconv.Atoi(request.Amount)
	events:=[]models.Event{}
	if a<=0{
		if err:=db.Database.Find(&events).Error;err!=nil{
			return e.ErrorFetching("events",c,err)
		}
	}
	if err:=db.Database.Limit(a).Find(&events).Error;err!=nil{
		return e.ErrorFetching("events",c,err)
	}
	return c.JSON(events)
}

func JoinToEvent(c *fiber.Ctx) error{
	eventId:=c.Params("id")
	var request struct{
		Id 		string 		`json:"user_id"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	user:=models.User{}
	// userId,_:=uuid.Parse(request.Id)
	if err := db.Database.First(&user,"id=?",request.Id).Error; err != nil {
		return e.NotFound("User",err,c)
    }
	event:=models.Event{}
	// eventIdUUID,_:=uuid.Parse(eventId)
	if err := db.Database.First(&event,"id=?",eventId).Error; err != nil {
        return e.NotFound("Event",err,c)
    }
	event.Members = append(event.Members, request.Id)
	user.Events = append(user.Events, eventId)
	if err:=db.Database.Save(&event).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Failed to join to event",
		})
	}

	if err:=db.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Failed to save event in user's events",
		})
	}

	return c.JSON(fiber.Map{
        "message": "Joined to event successfully",
    })
}

func UnjoinFromEvent(c *fiber.Ctx) error{
	userId:=c.Params("id")
	var request struct{
		Id 		string 		`json:"event_id"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	user:=models.User{}
	// userIdUUID,_:=uuid.Parse(userId)
	if err := db.Database.First(&user,"id=?",userId).Error; err != nil {
        return e.NotFound("User",err,c)
    }
	updateEvents:=make([]string,0,len(user.Events))
		for _, e := range user.Events {
			if e != request.Id {
				updateEvents = append(updateEvents, e)
				}
			}
	user.Events = updateEvents
	if err:=db.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"Failed to update user events",
		})
	}
	return c.JSON(fiber.Map{
		"message":"Unjoined from event successfully",
	})
}