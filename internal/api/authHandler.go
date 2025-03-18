package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/ctchen1999/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (a *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var params types.AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := a.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("Invalid credentials")
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return fmt.Errorf("Invalid Password")
	}

	resp := types.AuthResponse{
		User:  user,
		Token: GenerateToken(user),
	}

	fmt.Println("User logged in ->", user.FirstName, user.LastName)

	return response.SuccessResponse(c, resp)
}

func GenerateToken(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 1)
	claims := jwt.MapClaims{
		"id":      user.Id,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Failed to sign token with secret")
	}

	return tokenStr
}
