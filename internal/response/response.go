package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Code   int `json:"code"`
	Extras any `json:"extras"`
}

func NewResponse(code int, extras any) Response {
	return Response{
		Code:   code,
		Extras: extras,
	}
}

func SuccessResponse(c *fiber.Ctx, extras any) error {
	return c.JSON(NewResponse(http.StatusOK, fiber.Map{"data": extras}))
}

func ErrorResponse(c *fiber.Ctx, extras any) error {
	return c.JSON(NewResponse(http.StatusInternalServerError, extras))
}
