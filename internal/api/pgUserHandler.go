package api

import (
	models "github.com/ctchen1999/hotel-system/internal/pg"
	"github.com/ctchen1999/hotel-system/internal/pgtypes"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/gofiber/fiber/v2"
)

type PgUserHandler struct {
	userStore models.PgUserStore
}

func NewPgUserHandler(userStore models.PgUserStore) *PgUserHandler {
	return &PgUserHandler{
		userStore: userStore,
	}
}

func (h *PgUserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, users)
}

func (h *PgUserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, user)
}

func (h *PgUserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params pgtypes.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if validationErrors := params.Validate(); len(validationErrors) > 0 {
		return response.ErrorResponse(c, validationErrors)
	}

	user, err := pgtypes.NewUserFromParams(params)
	if err != nil {
		return err
	}

	if err := h.userStore.CreateUser(c.Context(), user); err != nil {
		return err
	}

	return response.SuccessResponse(c, user)
}

func (h *PgUserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		return err
	}

	return response.SuccessResponse(c, "User deleted successfully")
}

func (h *PgUserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	var params *pgtypes.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := h.userStore.UpdateUser(c.Context(), params, userId); err != nil {
		return err
	}

	return response.SuccessResponse(c, "User updated successfully")
}
