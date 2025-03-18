package api

import (
	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/ctchen1999/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	var params types.CreateHotelParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	hotel := &types.Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rating:   params.Rating,
	}
	createdHotel, err := h.store.Hotel.Create(c.Context(), hotel)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, createdHotel)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	//TODO: Implement query for hotels
	var query types.HotelQuery
	if err := c.QueryParser(&query); err != nil {
		return err
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelById(c.Context(), id)
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"hotelId": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, rooms)
}

func (h *HotelHandler) HandleUpdateHotel(c *fiber.Ctx) error {
	var params types.HotelUpdateParams
	hotelId := c.Params("id")

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := h.store.Hotel.Update(c.Context(), params, hotelId); err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.Map{"message": "hotel updated"})
}
