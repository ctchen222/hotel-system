package api

import (
	models "github.com/ctchen1999/hotel-system/internal/pg"
	"github.com/ctchen1999/hotel-system/internal/pgtypes"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/gofiber/fiber/v2"
)

type PgHotelHandler struct {
	hotelStore models.PgHotelStore
	roomStore  models.PgRoomStore
}

func NewPgHotelHandler(hotelStore models.PgHotelStore, roomStore models.PgRoomStore) *PgHotelHandler {
	return &PgHotelHandler{
		hotelStore: hotelStore,
		roomStore:  roomStore,
	}
}

func (h *PgHotelHandler) HandleCreateHotel(c *fiber.Ctx) error {
	var params pgtypes.CreateHotelParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	hotel := &pgtypes.Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rating:   params.Rating,
	}

	if err := h.hotelStore.CreateHotel(c.Context(), hotel); err != nil {
		return err
	}

	return response.SuccessResponse(c, "Hotel has been created.")
}

func (h *PgHotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, hotels)
}

func (h *PgHotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	hotelId := c.Params("id")
	hotel, err := h.hotelStore.GetHotelById(c.Context(), hotelId)
	if err != nil {
		return nil
	}

	return response.SuccessResponse(c, hotel)
}

func (h *PgHotelHandler) HandleUpdateHotel(c *fiber.Ctx) error {
	hotelId := c.Params("id")
	var params pgtypes.UpdateHotelParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := h.hotelStore.UpdateHotel(c.Context(), &params, hotelId); err != nil {
		return err
	}

	return response.SuccessResponse(c, "Hotel has been updated.")
}

func (h *PgHotelHandler) HandlerDeleteHotel(c *fiber.Ctx) error {
	hotelId := c.Params("id")
	if err := h.hotelStore.DeleteHotel(c.Context(), hotelId); err != nil {
		return err
	}
	return response.SuccessResponse(c, "Hotel has been deleted.")
}

func (h *PgHotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	hotelId := c.Params("id")

	rooms, err := h.roomStore.GetRooms(c.Context(), hotelId)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, rooms)
}
