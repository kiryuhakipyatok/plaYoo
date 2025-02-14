package errorpatterns

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func NotFound(what string,err error, c *fiber.Ctx) error{
	c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": fmt.Sprintf("%s not found: %s",what,err),
		})
}

func BadRequest(c *fiber.Ctx,err error) error{
	c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": fmt.Sprintf("Error parse request: %s",err),
	})
}

func BadUUID(c *fiber.Ctx, err error) error{
	c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": fmt.Sprintf("Error parse to uuid: %s",err),
	})
}

func ErrorFetching(what string,c *fiber.Ctx, err error) error{
	c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": fmt.Sprintf("Error fetching %s: %s",what,err),
	})
}