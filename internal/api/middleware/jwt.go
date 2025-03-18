package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get token from bearer
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return response.ErrUnAuthorized()
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := validateToken(token)
		if err != nil {
			return response.ErrUnAuthorized()
		}

		expires, err := time.Parse(time.RFC3339, claims["expires"].(string))
		if err != nil {
			return response.ErrUnAuthorized()
		}

		if time.Now().After(expires) {
			return response.ErrUnAuthorized()
		}

		userId := claims["id"]
		user, err := userStore.GetUserById(c.Context(), userId.(string))
		if err != nil {
			return response.ErrUnAuthorized()
		}

		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

// validate token
func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("UnAuthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		_ = fmt.Errorf("failed to parse jwt token = %s", tokenStr)
		return nil, fmt.Errorf("UnAuthorized")
	}

	// Payload
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("UnAuthorized")
}
