package api

import (
	"fmt"
	"time"

	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/ctchen1999/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var rawParams types.BookingRawParams
	if err := c.BodyParser(&rawParams); err != nil {
		return err
	}

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return response.ErrInvalidLocation()
	}

	from, err := time.Parse("2006-01-02", rawParams.From)
	if err != nil {
		return response.ErrInvalidDate()
	}
	from = from.In(loc)

	to, err := time.Parse("2006-01-02", rawParams.To)
	if err != nil {
		return response.ErrInvalidDate()
	}
	to = to.In(loc)

	params := types.BookingParams{
		From:      from,
		To:        to,
		NumPerson: rawParams.NumPerson,
	}
	if validationErrors := params.Validate(); len(validationErrors) > 0 {
		return response.ErrorResponse(c, validationErrors)
	}

	roomId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return response.ErrUnAuthenticated()
	}

	filter := bson.M{
		"roomId": roomId,
		"from": bson.M{
			"$gte": from,
		},
		"to": bson.M{
			"$lte": to,
		},
	}
	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return err
	}
	if len(bookings) > 0 {
		return response.ErrorResponse(c, fmt.Sprintf("Room %s is already booked", roomId.Hex()))
	}

	booking := types.Booking{
		UserId:    user.Id,
		RoomId:    roomId,
		NumPerson: params.NumPerson,
		From:      params.From,
		To:        params.To,
	}

	bookedRoom, err := h.store.Booking.InsertBookRoom(c.Context(), &booking)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, bookedRoom)
}

func (h *RoomHandler) HandleGetBookings(c *fiber.Ctx) error {
	var query types.BookingQuery
	if err := c.QueryParser(&query); err != nil {
		return err
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, bookings)
}
