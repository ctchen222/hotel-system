package api

import (
	"fmt"

	models "github.com/ctchen1999/hotel-system/internal/pg"
	"github.com/ctchen1999/hotel-system/internal/pgtypes"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/gofiber/fiber/v2"
)

type PgRoomHandler struct {
	roomStore models.PgRoomStore
}

func NewPgRoomHandler(roomStore models.PgRoomStore) *PgRoomHandler {
	return &PgRoomHandler{
		roomStore: roomStore,
	}
}

func (h *PgRoomHandler) HandlerGetRooms(c *fiber.Ctx) error {
	hotelId := c.Params("hotelId")
	fmt.Println("hotelId", hotelId)

	rooms, err := h.roomStore.GetRooms(c.Context(), hotelId)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, rooms)
}

func (h *PgRoomHandler) HandleCreateRoom(c *fiber.Ctx) error {
	hotelId := c.Params("hotelId")
	var params pgtypes.CreateRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := h.roomStore.CreateRoom(c.Context(), params, hotelId); err != nil {
		return err
	}

	return response.SuccessResponse(c, "Room has been created.")
}
