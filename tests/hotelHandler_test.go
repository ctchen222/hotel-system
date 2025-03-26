package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ctchen1999/hotel-system/internal/api"
	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/db/mocks"
	"github.com/ctchen1999/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

type HotelSuiteHandler struct {
	suite.Suite
	mockHotelStore *mocks.MockHotelStore
	mockRoomStore  *mocks.MockRoomStore
	hotelHandler   *api.HotelHandler

	hotels_1, hotels_2 []*types.Hotel
	hotel_embed        *types.HotelEmbed
	room_1, room_2     *types.Room
	rooms_1            []*types.Room
}

func (suite *HotelSuiteHandler) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.mockHotelStore = mocks.NewMockHotelStore(ctrl)
	suite.mockRoomStore = mocks.NewMockRoomStore(ctrl)
	mockUserStore := mocks.NewMockUserStore(ctrl)
	mockBookingStore := mocks.NewMockBookingStore(ctrl)
	store := &db.Store{
		User:    mockUserStore,
		Hotel:   suite.mockHotelStore,
		Room:    suite.mockRoomStore,
		Booking: mockBookingStore,
	}
	suite.hotelHandler = api.NewHotelHandler(store)
}

func (suite *HotelSuiteHandler) BeforeTest(suiteName, testName string) {
	hotelId1 := primitive.NewObjectID()
	hotelId2 := primitive.NewObjectID()
	hotelId3 := primitive.NewObjectID()
	hotelId4 := primitive.NewObjectID()
	hotelId5 := primitive.NewObjectID()
	suite.hotels_1 = []*types.Hotel{
		{
			Id:       hotelId1,
			Name:     "Hotel 1",
			Location: "Location 1",
			Rating:   5,
		},
		{
			Id:       hotelId2,
			Name:     "Hotel 2",
			Location: "Location 2",
			Rating:   4,
		},
	}
	suite.hotels_2 = []*types.Hotel{
		{
			Id:       hotelId3,
			Name:     "Hotel 3",
			Location: "Location 3",
			Rating:   5,
		},
		{
			Id:       hotelId4,
			Name:     "Hotel 4",
			Location: "Location 4",
			Rating:   4,
		},
		{
			Id:       hotelId5,
			Name:     "Hotel 5",
			Location: "Location 5",
			Rating:   5,
		},
	}

	roomId1 := primitive.NewObjectID()
	roomId2 := primitive.NewObjectID()
	suite.hotel_embed = &types.HotelEmbed{
		Id:       hotelId1,
		Name:     "Hotel 1",
		Location: "Location 1",
		Rating:   5,
		Rooms: []types.Room{
			{
				Id:      roomId1,
				Size:    "Single",
				SeaSide: true,
				Price:   100,
				HotelId: hotelId1,
			},
			{
				Id:      roomId2,
				Size:    "Double",
				SeaSide: false,
				Price:   200,
				HotelId: hotelId1,
			},
		},
	}

	suite.rooms_1 = []*types.Room{
		{
			Id:      roomId1,
			Size:    "Single",
			SeaSide: true,
			Price:   100,
			HotelId: hotelId1,
		},
		{
			Id:      roomId2,
			Size:    "Double",
			SeaSide: false,
			Price:   200,
			HotelId: hotelId1,
		},
	}

	// TODO: Setup app and routes here
}

func (suite *HotelSuiteHandler) TestHotelHandler_HandleGetHotels() {
	suite.mockHotelStore.EXPECT().GetHotels(gomock.Any(), gomock.Any()).Return(suite.hotels_1, nil).Times(1)
	suite.mockHotelStore.EXPECT().GetHotels(gomock.Any(), gomock.Any()).Return(suite.hotels_2, nil).Times(1)

	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		want    []*types.Hotel
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{c: &fiber.Ctx{}},
			want:    suite.hotels_1,
			wantErr: false,
		},
		{
			name:    "success",
			args:    args{c: &fiber.Ctx{}},
			want:    suite.hotels_2,
			wantErr: false,
		},
	}

	app := fiber.New()
	app.Get("/hotel", suite.hotelHandler.HandleGetHotels)

	type ResponseExtra struct {
		Data []*types.Hotel `json:"data"`
	}

	type r struct {
		Code   int           `json:"code"`
		Extras ResponseExtra `json:"extras"`
	}
	var response r

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/hotel", nil)
			resp, _ := app.Test(req)
			suite.Equal(http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			json.Unmarshal(body, &response)

			suite.Equal(len(tt.want), len(response.Extras.Data))
			suite.Equal(response.Code, http.StatusOK)
		})
	}
}

func (suite *HotelSuiteHandler) TestHotelHandler_HandleGetHotel() {
	suite.mockHotelStore.EXPECT().GetHotelById(gomock.Any(), gomock.Any()).Return(suite.hotel_embed, nil).Times(1)

	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		want    *types.HotelEmbed
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{c: &fiber.Ctx{}},
			want:    suite.hotel_embed,
			wantErr: false,
		},
	}

	app := fiber.New()
	app.Get("/hotel/:id", suite.hotelHandler.HandleGetHotel)

	type ResponseExtra struct {
		Data *types.HotelEmbed `json:"data"`
	}

	type r struct {
		Code   int           `json:"code"`
		Extras ResponseExtra `json:"extras"`
	}
	var response r

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/hotel/:id", nil)
			resp, _ := app.Test(req)
			suite.Equal(http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			json.Unmarshal(body, &response)

			suite.Equal(tt.want, response.Extras.Data)
			suite.Equal(response.Code, http.StatusOK)
		})
	}
}

func (suite *HotelSuiteHandler) TestHotelHandler_HandleGetRooms() {
	suite.mockRoomStore.EXPECT().GetRooms(gomock.Any(), gomock.Any()).Return(suite.rooms_1, nil).Times(1)
	paramsId := primitive.NewObjectID().Hex()

	type args struct {
		c        *fiber.Ctx
		paramsId string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "WrongParamsId",
			args: args{
				c:        &fiber.Ctx{},
				paramsId: "wrongParamsId",
			},
			want:    "the provided hex string is not a valid ObjectID",
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				c:        &fiber.Ctx{},
				paramsId: paramsId,
			},
			want:    suite.rooms_1,
			wantErr: false,
		},
	}

	app := fiber.New()
	app.Get("/hotel/:id/rooms", suite.hotelHandler.HandleGetRooms)

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/hotel/"+tt.args.paramsId+"/rooms", nil)
			resp, err := app.Test(req)
			body, _ := io.ReadAll(resp.Body)

			if (err != nil) != tt.wantErr {
				suite.Equal(tt.want, string(body))
				return
			}

			suite.Equal(http.StatusOK, resp.StatusCode)
		})
	}
}

func TestHotelSuiteHandler(t *testing.T) {
	suite.Run(t, new(HotelSuiteHandler))
}
