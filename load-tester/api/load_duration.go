package api

import (
	"context"

	"load-tester/service"
	"load-tester/util/errs"

	"github.com/gofiber/fiber/v2"
)

func (api *Api) loadDuration(c *fiber.Ctx) error {
	// Parse request configuration
	var param *service.LoadDurationParam
	if err := c.BodyParser(&param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"remark":         "failed to parse request body",
			"status_code":    errs.CODE_ERR_VALIDATION,
			"status_message": "failed to parse request body",
		})
	}

	// Call service layer
	result, err := api.service.LoadDuration(context.TODO(), param)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{
			"remark": err.Error(),
		})
	}

	return c.JSON(result)
}
