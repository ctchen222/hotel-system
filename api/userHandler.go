package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/ctchen1999/hotel-system/db"
	"github.com/ctchen1999/hotel-system/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)

	user, err := h.userStore.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no user found")
			return c.JSON(fiber.Map{"error": "no users found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	fmt.Println("TEST")
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	// validate the input
	if validationErrors := params.Validate(); len(validationErrors) > 0 {
		return c.JSON(fiber.Map{"errors": validationErrors})
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	createdUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(createdUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.userStore.DeleteUserById(c.Context(), userId); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "user deleted"})
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var params types.UserUpdateParams
	userId := c.Params("id")

	if err := c.BodyParser(&params); err != nil {
		return nil
	}

	if err := h.userStore.UpdateUser(c.Context(), params, userId); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "user updated"})
}
