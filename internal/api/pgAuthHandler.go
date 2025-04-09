package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	models "github.com/ctchen222/hotel-system/internal/pg"
	"github.com/ctchen222/hotel-system/internal/pgtypes"
	"github.com/ctchen222/hotel-system/internal/response"
	"github.com/ctchen222/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type PgAuthHandler struct {
	userStore models.PgUserStore
}

func NewPgAuthHandler(userStore models.PgUserStore) *PgAuthHandler {
	return &PgAuthHandler{
		userStore: userStore,
	}
}

func (a *PgAuthHandler) HandleLogin(c *fiber.Ctx) error {
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

	resp := pgtypes.PgAuthResponse{
		User:  user,
		Token: PgGenerateToken(user),
	}

	fmt.Println("User logged in ->", user.FirstName, user.LastName)

	return response.SuccessResponse(c, resp)
}

func PgGenerateToken(user *pgtypes.PGUser) string {
	now := time.Now()
	expires := now.Add(time.Hour * 24 * 7)
	claims := jwt.MapClaims{
		"id":      user.Id,
		"email":   user.Email,
		"expires": expires,
	}

	// header + payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	// sign token with secret -> hash_alg(header + payload + secret)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Failed to sign token with secret")
	}

	return tokenStr
}
