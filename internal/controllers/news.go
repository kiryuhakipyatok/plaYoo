package controllers

import (
	"avantura/backend/internal/db"
	"avantura/backend/internal/models"
	"time"
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
	"os"
	"fmt"
	"path/filepath"
	e "avantura/backend/pkg/error-patterns"
	"strconv"
)


func AddNews(c *fiber.Ctx) error{
	var request struct{
		Title 		string 		`json:"title"`
		Body 		string 		`json:"body"`
		Link    	string 		`json:"link"`
		AuthorName  string		`json:"author_name"`
	}

	request.Title = c.FormValue("title")
    request.Body = c.FormValue("body")
    request.Link = c.FormValue("link")
    request.AuthorName = c.FormValue("author_name")

	file, err := c.FormFile("picture")
    if err != nil {
        c.Status(fiber.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }

	authorId:=c.Params("id")
	// authorIdUUID,_:=uuid.Parse(authorId)
	author:=models.User{}
	if err:=db.Database.First(&author,"id=?",authorId).Error;err!=nil{
		return e.NotFound("Author",err,c)
	}

	news:=models.News{
		Id:authorId+"news",
		AuthorName: request.AuthorName,
		AuthorId: authorId,
		AuthorAvatars: author.Avatar,
		Title: request.Title,
		Body: request.Body,
		Time:time.Now(),
		Link: request.Link,
	}
	
	uploadDir := "../../pkg/news_picture"
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	fileName := fmt.Sprintf("%s-news-picture%s", news.Id,filepath.Ext(file.Filename))
    filepath := filepath.Join(uploadDir, fileName)
	if _, err := os.Stat(news.Picture); err == nil {
		if err := os.Remove(news.Picture); err != nil {
			c.Status(fiber.StatusInternalServerError)
        	return c.JSON(fiber.Map{
            "error": fmt.Sprintf("Failed to remove file: %v", err),
        })
		}
	}
    if err := c.SaveFile(file, filepath); err != nil {
        c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": fmt.Sprintf("Failed to save file: %v", err),
        })
    }
	fileURL:=fmt.Sprintf("http://localhost:9110/pkg/news_picture/%s",fileName)
	news.Picture = fileURL
	if err:=db.Database.Create(&news).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error creating news",
		})
	}

	return c.JSON(news)
}


func GetNews(c *fiber.Ctx) error{
	// var request struct{
	// 	Amount string `json:"amount"`
	// }
	amount := c.Query("amount")
	// if err:=c.BodyParser(&request);err!=nil{
	// 	return e.BadRequest(c,err)
	// }
	amountI,_:=strconv.Atoi(amount)
	//a,_:=strconv.Atoi(request.Amount)
	news:=[]models.News{}
	// if err:=db.Database.Limit(a).Find(&news).Error;err!=nil{
	// 	return e.ErrorFetching("news",c,err)
	// }
	if err:=db.Database.Limit(amountI).Find(&news).Error;err!=nil{
		return e.ErrorFetching("news",c,err)
	}
	return c.JSON(news)
}


func GetConcreteNews(c *fiber.Ctx) error{
	id:=c.Params("id")
	news:=models.News{}
	if err:=db.Database.Find(&news,"id=?",id).Error;err!=nil{
		return e.ErrorFetching("news",c,err)
	}
	return c.JSON(news)
}
