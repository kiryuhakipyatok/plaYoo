package controllers

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	e "avantura/backend/storage/error-patterns"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddComment(c *fiber.Ctx) error{
	commentdata:=map[string]string{}
	if err:=c.BodyParser(&commentdata);err!=nil{
		return e.BadRequest(c,err)
	}
	authorId:=commentdata["author_id"]	

	author:=models.User{}
	if err:=postgres.Database.First(&author,"id=?",authorId).Error;err!=nil{
		return e.NotFound("Author",err,c)
	}
	authorIdUUID,err:=uuid.Parse(authorId)
	if err != nil {
		return e.BadUUID(c,err)
	}
	var tempId uuid.UUID
	comment:=models.Comment{
		Id:uuid.New(),
		AuthorId: authorIdUUID,
		AuthorName: author.Login,
		AuthorAvatar: author.Avatar,
		Body: commentdata["body"],
		Time: time.Now(),
		Receiver: tempId,
	} 
	if userId,ok:=commentdata["user_id"];ok && userId!=""{
		id,err:=uuid.Parse(userId)
		if err != nil {
			return e.BadUUID(c,err)
		}
		comment.Receiver = id
		user:=models.User{}
		if err:=postgres.Database.First(&user,"id=?",userId).Error;err!=nil{
			return e.NotFound("User",err,c)
		}

		user.Comments = append(user.Comments, comment.Id.String())
		postgres.Database.Save(&user)
	}else if eventId,ok:=commentdata["event_id"];ok && eventId!=""{
		id,err:=uuid.Parse(eventId)
		if err != nil {
			return e.BadUUID(c,err)
		}
		comment.Receiver = id
		event:=models.Event{}
		if err:=postgres.Database.First(&event,"id=?",eventId).Error;err!=nil{
			return e.NotFound("Event",err,c)
		}

		event.Comments = append(event.Comments, comment.Id.String())
		postgres.Database.Save(&event)
	}else if newsId,ok:=commentdata["news_id"];ok && newsId!=""{
		id,err:=uuid.Parse(newsId)
		if err != nil {
			return e.BadUUID(c,err)
		}
		comment.Receiver = id
		news:=models.News{}
		if err:=postgres.Database.First(&news,"id=?",newsId).Error;err!=nil{
			return e.NotFound("News",err,c)
		}

		news.Comments = append(news.Comments, comment.Id.String())
		postgres.Database.Save(&news)
	}else{
		return e.NotFound("Receiver",err,c)
	}
	
	
	if err:=postgres.Database.Create(&comment).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error creating comment",
		})
	}

	return c.JSON(comment)

}


func ShowComments(c *fiber.Ctx) error{
	id :=c.Query("id")
	amount:=c.Query("amount")
	if amount=="" && id ==""{
		var request struct{
			UserId string `json:"user_id"`
			Amount string `json:"amount"`
		}
		if err:=c.BodyParser(&request);err!=nil{
			return e.BadRequest(c,err)
		}
		amount = request.Amount
		id = request.UserId
	}
	amountI,_:=strconv.Atoi(amount)
	user:=models.User{}
	if err:=postgres.Database.First(&user,"id=?",id).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	amountComments := user.Comments[:amountI]
	return c.JSON(amountComments)
}