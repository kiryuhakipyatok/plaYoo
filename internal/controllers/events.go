package controllers

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"avantura/backend/internal/db/redis"
	e "avantura/backend/pkg/error-patterns"
	"encoding/json"
	"log"
	"time"

	r "github.com/redis/go-redis/v9"
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
	if err:=postgres.Database.First(&user,"id=?",author_id).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	max,_:=strconv.Atoi(iventdata["max"])
	minute,_:=strconv.Atoi(iventdata["minute"])
	members:=[]string{author_id}
	event:=models.Event{
		Id:	uuid.New(),
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
	game:=models.Game{}
	if err:=postgres.Database.Transaction(func(tx *gorm.DB) error {
		if err:=tx.Create(&event).Error;err!=nil{
			return err
		}
		eventData,err:=json.Marshal(&event)
		if err!=nil{
			return err
		}
		ttl:=time.Minute*time.Duration(minute)
		if redis.Rdb!=nil{
			if err:=redis.Rdb.Set(redis.Ctx,event.Id.String(),eventData,ttl).Err();err!=nil{
				return err
			}
		}else{
			log.Printf("Redis is not connected")
		}
		if err:=tx.First(&game,"name=?",iventdata["game"]).Error; err != nil {
			return err
		}
		user.Events=append(user.Events, event.Id.String())
		game.NumberOfEvents++
		tx.Save(&user)
		tx.Save(&game)
		return nil
	});err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error transaction of creating event",
		})
	}
	// if err:=postgres.Database.Create(&event).Error;err!=nil{
	// 	c.Status(fiber.StatusInternalServerError)
	// 	return c.JSON(fiber.Map{
	// 		"error": "Error creating ivent",
	// 	})
	// }
	
	if err := postgres.Database.First(&game,"name=?",iventdata["game"]).Error; err != nil {
        return e.NotFound("Game",err,c)
    }
	user.Events=append(user.Events, event.Id.String())
	game.NumberOfEvents++
	postgres.Database.Save(&user)
	postgres.Database.Save(&game)
	return c.JSON(event)
}

func GetConcreteEvent(c *fiber.Ctx) error{
	eventId:=c.Params("id")
	event:=models.Event{}
	data,err:=redis.Rdb.Get(redis.Ctx,eventId).Result()
	if err!=nil{
		if err==r.Nil{
			eventData, _ := json.Marshal(data)
			ttl := time.Until(event.Time)
            redis.Rdb.Set(redis.Ctx, eventId, eventData, ttl)
			if err:=postgres.Database.First(&event,"id=?",eventId).Error;err!=nil{
				return e.NotFound("User",err,c)
			}
		}else{
			if err:=postgres.Database.First(&event,"id=?",eventId).Error;err!=nil{
				return e.NotFound("Event",err,c)
			}
			log.Printf("Error getting event from Redis, getting event from Postgre")
		}
	}else{
		if err:=json.Unmarshal([]byte(data),&event);err!=nil{
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"error": "Failed unmarshal",
			})
		}
	}
	return c.JSON(event)
}

func GetEvents(c *fiber.Ctx) error{
	amount:=c.Query("amount")
	if amount==""{
		var request struct{
			Amount string `json:"amount"`
		}
		if err:=c.BodyParser(&request);err!=nil{
			return e.BadRequest(c,err)
		}
		amount = request.Amount
	}
	amountI,_:=strconv.Atoi(amount)
	events:=[]models.Event{}
	if err:=postgres.Database.Limit(amountI).Find(&events).Error;err!=nil{
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
	if err := postgres.Database.First(&user,"id=?",request.Id).Error; err != nil {
		return e.NotFound("User",err,c)
    }
	event:=models.Event{}
	// eventIdUUID,_:=uuid.Parse(eventId)
	if err := postgres.Database.First(&event,"id=?",eventId).Error; err != nil {
        return e.NotFound("Event",err,c)
    }
	event.Members = append(event.Members, request.Id)
	user.Events = append(user.Events, eventId)
	if err:=postgres.Database.Save(&event).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Failed to join to event",
		})
	}

	if err:=postgres.Database.Save(&user).Error;err!=nil{
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
	if err := postgres.Database.First(&user,"id=?",userId).Error; err != nil {
        return e.NotFound("User",err,c)
    }
	updateEvents:=make([]string,0,len(user.Events))
		for _, e := range user.Events {
			if e != request.Id {
				updateEvents = append(updateEvents, e)
				}
			}
	user.Events = updateEvents
	if err:=postgres.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"Failed to update user events",
		})
	}
	return c.JSON(fiber.Map{
		"message":"Unjoined from event successfully",
	})
}