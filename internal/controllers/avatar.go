package controllers

import(
	"avantura/backend/internal/db"
	"avantura/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"os"
	"fmt"
	"path/filepath"
	e "avantura/backend/pkg/error-patterns"
)

func UploadAvatar(c *fiber.Ctx) error{
	userId := c.Params("id")

    file, err := c.FormFile("avatar")
    if err != nil {
        c.Status(fiber.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }

    uploadDir := "../../pkg/avatars"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}
	
    
	user := models.User{}
	if err:=db.Database.First(&user,"id=?",userId).Error;err!=nil{
		return e.NotFound("User",err,c)
	}
    fileName := fmt.Sprintf("%s-avatar%s", user.Login,filepath.Ext(file.Filename))
    filepath := filepath.Join(uploadDir, fileName)
	if _, err := os.Stat(user.Avatar); err == nil {
		if err := os.Remove(user.Avatar); err != nil {
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

	fileURL:=fmt.Sprintf("http://localhost:9110/pkg/avatars/%s",fileName)

	user.Avatar = fileURL
	if err := db.Database.Save(&user).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Failed to save user",
        })
	}

	return c.JSON(fiber.Map{
        "message": "Avatar uploaded successfully",
        "path":    fileURL,
    })
}

func DeleteAvatar(c *fiber.Ctx) error{
	userId:=c.Params("id")
	userIdUUID,_:=uuid.Parse(userId)
	user:=models.User{}
	if err:=db.Database.First(&user,"id=?",userIdUUID).Error;err!=nil{
		return e.NotFound("User",err,c)	
	}
	if user.Avatar == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "User does not have an avatar",
		})
	}
	if err := os.Remove(user.Avatar); err != nil {
		c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
        "error": fmt.Sprintf("Failed to remove file: %v", err),
		})
	}
	user.Avatar = ""
	db.Database.Save(&user)
	return c.JSON(fiber.Map{
        "message": "Avatar deleted successfully",
    })
}