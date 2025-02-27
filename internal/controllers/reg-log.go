package controllers

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	"time"
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"avantura/backend/pkg/constants"
	e "avantura/backend/pkg/error-patterns"
)

func Register(c *fiber.Ctx) error{
	userdata:=map[string]string{}
	if err:=c.BodyParser(&userdata);err!=nil{
		return e.BadRequest(c,err)
	}
	password,err:=bcrypt.GenerateFromPassword([]byte(userdata["password"]),14)
	if err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Error brcypt hash password",
		})
	}
	user:=models.User{
		Id:userdata["login"]+"id",
		Login: userdata["login"],
		Tg: userdata["tg"],
		Password: password,
	}
	if err:=postgres.Database.Create(&user).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Error creating user" + err.Error(),
		})
	}

	return c.JSON(user)

}

func Login(c *fiber.Ctx) error{

	userdata:=map[string]string{}
	if err:=c.BodyParser(&userdata);err!=nil{
		return e.BadRequest(c,err)
	}
	
	user:=models.User{}
	if err:=postgres.Database.First(&user,"login=?",userdata["login"]).Error;err!=nil{
		return e.NotFound("User",err,c)
	}

	if err:=bcrypt.CompareHashAndPassword(user.Password,[]byte(userdata["password"]));err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":"Incorrect password",
		})
	}

	claims:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
		Issuer: user.Id,
		ExpiresAt: time.Now().Add(time.Hour*24).Unix(),
	})

	token,err:=claims.SignedString([]byte(constants.Secret))
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Could not login",
		})
	}

	cookie:=fiber.Cookie{
		Name:"jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message":"Success",
	})
}

func Logout(c *fiber.Ctx) error{
	cookie:=fiber.Cookie{
		Name:"jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"Succes",
	})
}



