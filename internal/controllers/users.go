package controllers

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	"avantura/backend/storage/constants"
	e "avantura/backend/storage/error-patterns"
	"strconv"
	"avantura/backend/internal/db/redis"
	r "github.com/redis/go-redis/v9"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"encoding/json"
	"time"
	"log"
)

func GetAllUsers(c *fiber.Ctx) error{
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
	users:=[]models.User{}
	if err:=postgres.Database.Limit(amountI).Find(&users).Error;err!=nil{
		return e.ErrorFetching("users",c,err)
	}
	return c.JSON(users)
}

func GetConcreteUser(c *fiber.Ctx) error{
	searchUserId:=c.Params("id")
	user:=models.User{}
	data,err:=redis.Rdb.Get(redis.Ctx,searchUserId).Result()
	if err!=nil{
		if err==r.Nil{
			if err := postgres.Database.First(&user,"id=?",searchUserId).Error; err != nil {
				return e.NotFound("User",err,c)
			}
			userData, _ := json.Marshal(user)
			ttl:=time.Hour*24
            redis.Rdb.Set(redis.Ctx, searchUserId, userData, ttl)
		}else{
			if err := postgres.Database.First(&user,"id=?",searchUserId).Error; err != nil {
				return e.NotFound("User",err,c)
			}
			log.Printf("Error getting user from Redis, getting user from Postgre")
		}
	}else{
		if err:=json.Unmarshal([]byte(data),&user);err!=nil{
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"error": "Failed unmarshal",
			})
		}
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