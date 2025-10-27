package api

import (
	"context"

	"load-tester/service"
	"load-tester/util/errs"

	"github.com/gofiber/fiber/v2"
)

func (api *Api) loadRps(c *fiber.Ctx) error {
	// Parse request configuration
	var param service.LoadRpsParam
	if err := c.BodyParser(&param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"remark":         "failed to parse request body",
			"status_code":    errs.CODE_ERR_VALIDATION,
			"status_message": err.Error(),
		})
	}

	// Validate parameters
	if err := api.validateLoadRpsParam(&param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"remark":         "validation failed",
			"status_code":    errs.CODE_ERR_VALIDATION,
			"status_message": err.Error(),
		})
	}

	// Call service layer
	result, err := api.service.LoadRps(context.Background(), &param)
	if err != nil {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*errs.Error); ok && e.Kind == errs.Validation {
			code = fiber.StatusBadRequest
		}
		return c.Status(code).JSON(map[string]any{
			"remark":         "operation failed",
			"status_code":    errs.CODE_ERR_UNANTICIPATED,
			"status_message": err.Error(),
		})
	}

	return c.JSON(result)
}
