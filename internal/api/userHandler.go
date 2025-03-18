package api

import (
	"errors"
	"fmt"

	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/ctchen1999/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)

	user, err := h.store.User.GetUserById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no user found")
			return c.JSON(fiber.Map{"error": "no users found"})
		}
		return err
	}

	return response.SuccessResponse(c, user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.User.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	// validate the input
	if validationErrors := params.Validate(); len(validationErrors) > 0 {
		return response.ErrorResponse(c, validationErrors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	createdUser, err := h.store.User.Create(c.Context(), user)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, createdUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.store.User.DeleteById(c.Context(), userId); err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.Map{"message": "user deleted"})
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var params types.UserUpdateParams
	userId := c.Params("id")

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := h.store.User.Update(c.Context(), params, userId); err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.Map{"message": "user updated"})
}
