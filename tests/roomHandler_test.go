package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ctchen1999/hotel-system/internal/api"
	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/db/mocks"
	"github.com/ctchen1999/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

type RoomSuiteHandler struct {
	suite.Suite
	mockBookingStore *mocks.MockBookingStore
	roomHandler      *api.RoomHandler

	bookings []*types.Booking
}

func (suite *RoomSuiteHandler) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.mockBookingStore = mocks.NewMockBookingStore(ctrl)
	mockUserStore := mocks.NewMockUserStore(ctrl)
	mockHotelStore := mocks.NewMockHotelStore(ctrl)
	mockRoomStore := mocks.NewMockRoomStore(ctrl)
	store := &db.Store{
		User:    mockUserStore,
		Hotel:   mockHotelStore,
		Room:    mockRoomStore,
		Booking: suite.mockBookingStore,
	}
	suite.roomHandler = api.NewRoomHandler(store)
}

func (suite *RoomSuiteHandler) BeforeTest(suiteName, testName string) {
	suite.bookings = []*types.Booking{
		{
			Id:        primitive.NewObjectID(),
			UserId:    primitive.NewObjectID(),
			RoomId:    primitive.NewObjectID(),
			NumPerson: 1,
			From:      time.Now(),
			To:        time.Now().AddDate(0, 0, 1),
		},
		{
			Id:        primitive.NewObjectID(),
			UserId:    primitive.NewObjectID(),
			RoomId:    primitive.NewObjectID(),
			NumPerson: 3,
			From:      time.Now(),
			To:        time.Now().AddDate(0, 0, 3),
		},
	}
}

func (suite *RoomSuiteHandler) TestRoomHandler_HandleBookRoom() {
}

func (suite *RoomSuiteHandler) TestRoomHandler_HandleGetBookings() {
	suite.mockBookingStore.EXPECT().GetBookings(gomock.Any(), gomock.Any()).Return(
		suite.bookings, nil).AnyTimes()

	app := fiber.New()
	app.Get("/room/booking", suite.roomHandler.HandleGetBookings)

	req := httptest.NewRequest(http.MethodGet, "/room/booking", nil)
	resp, _ := app.Test(req)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func TestRoomSuiteHandler(t *testing.T) {
	suite.Run(t, new(RoomSuiteHandler))
}
