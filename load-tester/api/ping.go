package api

import (
	"context"

	"load-tester/service"

	"github.com/gofiber/fiber/v2"
)

func (api *Api) Ping(c *fiber.Ctx) error {
	// Parse request queries
	queries := c.Queries()

	// Call service layer
	result, err := api.service.Ping(context.TODO(), &service.PingParam{
		Message: queries["message"],
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"remark": err.Error(),
		})
	}

	return c.JSON(result)
}
