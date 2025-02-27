package controllers

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	"avantura/backend/pkg/constants"
	e "avantura/backend/pkg/error-patterns"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func GetAllUsers(c *fiber.Ctx) error{
	var request struct{
		Amount string `json:"amount"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	a,_:=strconv.Atoi(request.Amount)
	users:=[]models.User{}
	if err:=postgres.Database.Limit(a).Find(&users).Error;err!=nil{
		return e.ErrorFetching("users",c,err)
	}
	return c.JSON(users)
}

func GetConcreteUser(c *fiber.Ctx) error{
	searchUserId:=c.Params("id")
	user:=models.User{}
	if err := postgres.Database.First(&user,"id=?",searchUserId).Error; err != nil {
       return e.NotFound("User",err,c)
    }	
	return c.JSON(user)
}

func User(c *fiber.Ctx) error{
	cookie:=c.Cookies("jwt")
	token,err:=jwt.ParseWithClaims(cookie,&jwt.StandardClaims{},func(t *jwt.Token) (interface{},error){
		return []byte(constants.Secret),nil
	})
	
	if err!=nil{
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error":"Unauthenticated",
		})
	}
	claims:=token.Claims.(*jwt.StandardClaims)
	user:=models.User{}
	postgres.Database.Where("id=?",claims.Issuer).First(&user)
	return c.JSON(user)
}