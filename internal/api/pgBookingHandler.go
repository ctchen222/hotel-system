package api

import (
	"strconv"
	"time"

	models "github.com/ctchen1999/hotel-system/internal/pg"
	"github.com/ctchen1999/hotel-system/internal/pgtypes"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/gofiber/fiber/v2"
)

type PgBookingHandler struct {
	bookingStore models.BookingStore
}

func NewPgBookingHandler(bookingStore models.BookingStore) *PgBookingHandler {
	return &PgBookingHandler{
		bookingStore: bookingStore,
	}
}

func (h *PgBookingHandler) HandleCreateBooking(c *fiber.Ctx) error {
	var params pgtypes.BookingParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return response.ErrInvalidLocation()
	}

	from, err := time.Parse("2006-01-02", params.FromDate)
	if err != nil {
		return response.ErrInvalidDate()
	}
	from = from.In(loc)

	to, err := time.Parse("2006-01-02", params.ToDate)
	if err != nil {
		return response.ErrInvalidDate()
	}
	to = to.In(loc)

	roomId, err := strconv.Atoi(params.RoomId)
	if err != nil {
		return response.ErrParseInt()
	}

	user, ok := c.Context().UserValue("user").(*pgtypes.PGUser)
	if !ok {
		return response.ErrUnAuthenticated()
	}
	userId, err := strconv.Atoi(user.Id)
	if err != nil {
		return response.ErrParseInt()
	}

	bookingParams := pgtypes.Booking{
		UserId:    userId,
		RoomId:    roomId,
		FromDate:  from,
		ToDate:    to,
		NumPerson: params.NumPerson,
	}

	if err := h.bookingStore.CreateBooking(c.Context(), &bookingParams); err != nil {
		return err
	}

	return response.SuccessResponse(c, "Booking has been created.")
}

func (h *PgBookingHandler) HandleGetBookingInfo(c *fiber.Ctx) error {
	userId := c.Params("userId")

	bookingInfos, err := h.bookingStore.GetBookingByUserId(c.Context(), userId)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, bookingInfos)
}
