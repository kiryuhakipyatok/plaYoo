package controllers

import (
	"avantura/backend/internal/db/postgres"
	"avantura/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	e "avantura/backend/pkg/error-patterns"
)

func Follow(c *fiber.Ctx) error{
	var request struct{
		UserId		string	`json:"user_id"`
		FollowLogin string  `json:"follow_login"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	user:=models.User{}
	if err := postgres.Database.First(&user,"id=?",request.UserId).Error; err != nil {
		return e.NotFound("User",err,c)
    }
	if user.Login == request.FollowLogin{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Its you login",
		})
	}	

	follow:=models.User{}

	if err := postgres.Database.First(&follow,"login=?",request.FollowLogin).Error; err != nil {
        return e.NotFound("Follow",err,c)
    }

	for _,followId:=range user.Followings{
		if followId == follow.Id{
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message":"You alredy follow to this user",
			})
		}
	}

	user.Followings = append(user.Followings, follow.Id)
	follow.Followers = append(follow.Followers, user.Id)

	if err:=postgres.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Failed to follow user",
		})
	}

	if err:=postgres.Database.Save(&follow).Error;err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"Failed to add user to followers",
		})
	}

	return c.JSON(fiber.Map{
        "message": "success",
    })
}

func Unfollow(c *fiber.Ctx) error{
	var request struct{
		UserId		string		`json:"user_id"`
		FollowId 	string 		`json:"follow_id"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.BadRequest(c,err)
	}
	user:=models.User{}
	if err := postgres.Database.First(&user,"id=?",request.UserId).Error; err != nil {
		return e.NotFound("User",err,c)
    }
	follow:=models.User{}
	if err := postgres.Database.First(&follow,"id=?",request.FollowId).Error; err != nil {
		return e.NotFound("Follow",err,c)
    }
	updateFollowings:=make([]string,0,len(user.Followings))
	for _, f := range user.Followings {
		if f != request.FollowId {
			updateFollowings = append(updateFollowings, f)
		}
	}
	updateFollowers:=make([]string,0,len(follow.Followers))
	for _, f := range user.Followers {
		if f != request.FollowId {
			updateFollowers = append(updateFollowers, f)
		}
	}
	if len(updateFollowings) != 0{
		user.Followings = updateFollowings
	}else{
		user.Followings = nil
	}
	if len(updateFollowers)!=0{
		follow.Followers = updateFollowers
	}else{
		follow.Followers = nil
	}
	
	
	if err:=postgres.Database.Save(&user).Error;err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"Failed to update user followings",
		})
	}
	if err:=postgres.Database.Save(&follow).Error;err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"Failed to update follow followers",
		})
	}
	return c.JSON(fiber.Map{
        "message": "success",
    })
}