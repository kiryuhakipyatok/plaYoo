package controllers

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	e "avantura/backend/pkg/error-patterns"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
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
	// authorIdUUID,err:=uuid.Parse(authorId)
	// if err != nil {
	// 	return e.BadUUID(c,err)
	// }
	var tempId string
	comment:=models.Comment{
		Id:authorId+"comment"+tempId,
		AuthorId: authorId,
		AuthorName: author.Login,
		AuthorAvatar: author.Avatar,
		Body: commentdata["body"],
		Time: time.Now(),
		Receiver: tempId,
	} 
	if userId,ok:=commentdata["user_id"];ok && userId!=""{
		// id,err:=uuid.Parse(userId)
		// if err != nil {
		// 	return e.BadUUID(c,err)
		// }
		comment.Receiver = userId
		user:=models.User{}
		if err:=postgres.Database.First(&user,"id=?",userId).Error;err!=nil{
			return e.NotFound("User",err,c)
		}

		user.Comments = append(user.Comments, comment.Id)
		postgres.Database.Save(&user)
	}else if eventId,ok:=commentdata["event_id"];ok && eventId!=""{
		// id,_:=uuid.Parse(eventId)
		comment.Receiver = eventId
		event:=models.Event{}
		if err:=postgres.Database.First(&event,"id=?",eventId).Error;err!=nil{
			return e.NotFound("Event",err,c)
		}

		event.Comments = append(event.Comments, comment.Id)
		postgres.Database.Save(&event)
	}else if newsId,ok:=commentdata["news_id"];ok && newsId!=""{
		// id,_:=uuid.Parse(newsId)
		comment.Receiver = newsId
		news:=models.News{}
		if err:=postgres.Database.First(&news,"id=?",newsId).Error;err!=nil{
			return e.NotFound("News",err,c)
		}

		news.Comments = append(news.Comments, comment.Id)
		postgres.Database.Save(&news)
	}else{
		// return e.NotFound("Receiver",err,c)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Receiver not found",
		})
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
	id:=c.Params("id")
	var request struct{
		Amount string `json:"amount"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	user:=models.User{}
	if err:=postgres.Database.First(&user,"id=?",id).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
	a,_:=strconv.Atoi(request.Amount)
	amountComments := user.Comments[:a]
	return c.JSON(amountComments)
}