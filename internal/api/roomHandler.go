package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ctchen1999/hotel-system/internal/db"
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
	var rawParams types.BooingRawParams
	if err := json.Unmarshal(c.BodyRaw(), &rawParams); err != nil {
		return c.JSON(fiber.Map{"error": "Invalid request body"})
	}

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return c.JSON(fiber.Map{"error": "Invalid location"})
	}

	from, err := time.Parse("2006-01-02", rawParams.From)
	if err != nil {
		return c.JSON(fiber.Map{"error": "Invalid from date"})
	}
	from = from.In(loc)

	to, err := time.Parse("2006-01-02", rawParams.To)
	if err != nil {
		return c.JSON(fiber.Map{"error": "Invalid to date"})
	}
	to = to.In(loc)

	params := types.BookingParams{
		From:      from,
		To:        to,
		NumPerson: rawParams.NumPerson,
	}
	if validationErrors := params.Validate(); len(validationErrors) > 0 {
		return c.JSON(fiber.Map{"errors": validationErrors})
	}

	roomId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
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
	fmt.Print("bookings", bookings)
	if err != nil {
		return err
	}
	if len(bookings) > 0 {
		return c.JSON(fiber.Map{
			"Type": "error",
			"Msg":  fmt.Sprintf("room %s has been booked", roomId),
		})
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

	return c.JSON(bookedRoom)
}

func (h *RoomHandler) HandleGetBookings(c *fiber.Ctx) error {
	var query types.BookingQuery
	if err := c.QueryParser(&query); err != nil {
		return err
	}
	fmt.Println("query", query)

	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}
