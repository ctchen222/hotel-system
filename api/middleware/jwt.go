package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("JWT Middleware")

	// token, ok := c.GetReqHeaders()["X-Api-Token"]
	// if !ok {
	// 	return fmt.Errorf("unAuthorized")
	// }

	// fmt.Println("Token: ", token)

	return c.Next()
}
