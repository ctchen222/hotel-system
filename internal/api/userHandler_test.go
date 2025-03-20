package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/ctchen1999/hotel-system/internal/api"
	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/db/mocks"
	"github.com/ctchen1999/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

type UserSuiteHandler struct {
	suite.Suite
	mockUserStore    *mocks.MockUserStore
	mockBookingStore *mocks.MockBookingStore
	mockHotelStore   *mocks.MockHotelStore
	mockRoomStore    *mocks.MockRoomStore
	userHandler      *api.UserHandler
}

func (suite *UserSuiteHandler) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.mockUserStore = mocks.NewMockUserStore(ctrl)
	suite.mockBookingStore = mocks.NewMockBookingStore(ctrl)
	suite.mockHotelStore = mocks.NewMockHotelStore(ctrl)
	suite.mockRoomStore = mocks.NewMockRoomStore(ctrl)

	store := &db.Store{
		User:    suite.mockUserStore,
		Booking: suite.mockBookingStore,
		Hotel:   suite.mockHotelStore,
		Room:    suite.mockRoomStore,
	}
	suite.userHandler = api.NewUserHandler(store)
}

func (suite *UserSuiteHandler) TestUserHandler_HandleGetUser() {
	userId := primitive.NewObjectID()
	user := &types.User{
		Id:                userId,
		FirstName:         "TwoBao",
		LastName:          "Chen",
		Email:             "twobao@twobao.com",
		EncryptedPassword: "123456",
	}

	suite.mockUserStore.EXPECT().GetUserById(gomock.Any(), userId.Hex()).Return(user, nil)

	app := fiber.New()
	app.Get("/users/:id", suite.userHandler.HandleGetUser)

	req := httptest.NewRequest("GET", "/users/"+userId.Hex(), nil)
	resp, err := app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)
}

func (suite *UserSuiteHandler) TestUserHandler_HandleGetUsers() {

	userId1 := primitive.NewObjectID()
	userId2 := primitive.NewObjectID()
	userId3 := primitive.NewObjectID()
	users := []*types.User{
		{
			Id:                userId1,
			FirstName:         "TwoBao",
			LastName:          "Chen",
			Email:             "twobao@twobao.com",
			EncryptedPassword: "123456",
		},
		{
			Id:                userId2,
			FirstName:         "Shaun",
			LastName:          "Chen",
			Email:             "shaun@twobao.com",
			EncryptedPassword: "123456",
		},
		{
			Id:                userId3,
			FirstName:         "Joanne",
			LastName:          "Lin",
			Email:             "joanne@twobao.com",
			EncryptedPassword: "123456",
		},
	}

	suite.mockUserStore.EXPECT().GetUsers(gomock.Any()).Return(
		users, nil)

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	// var resp response.Response
	err := suite.userHandler.HandleGetUsers(ctx)

	assert.NoError(suite.T(), err)
}

func (suite *UserSuiteHandler) TestUserHandler_HandlePostUser() {
	userId := primitive.NewObjectID()
	user := &types.User{
		Id:                userId,
		FirstName:         "TwoBao",
		LastName:          "Chen",
		Email:             "twobao@twobao.com",
		EncryptedPassword: "test12345678",
	}
	userCreateParams := &types.CreateUserParams{
		FirstName: "TwoBao",
		LastName:  "Chen",
		Email:     "twobao@twobao.com",
		Password:  "test1234",
	}
	userBody, _ := json.Marshal(userCreateParams)

	suite.mockUserStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(user, nil)

	app := fiber.New()
	app.Post("/users", suite.userHandler.HandlePostUser)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(userBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(suite.T(), 200, resp.StatusCode)
	assert.Equal(suite.T(), "application/json", resp.Header.Get("Content-Type"))
}

func (suite *UserSuiteHandler) TestUserHandler_HandleDeleteUser() {
	suite.mockUserStore.EXPECT().DeleteById(gomock.Any(), gomock.Any()).Return(nil)

	userId := primitive.NewObjectID()
	app := fiber.New()
	app.Delete("/users/:id", suite.userHandler.HandleDeleteUser)

	req := httptest.NewRequest("DELETE", "/users/"+userId.Hex(), nil)
	resp, _ := app.Test(req)
	assert.Equal(suite.T(), 200, resp.StatusCode)
}

func (suite *UserSuiteHandler) TestUserHandler_HandleUpdateUser() {
	suite.mockUserStore.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	userId := primitive.NewObjectID()
	userUpdateParams := &types.UserUpdateParams{
		FirstName: "TwoBao",
		LastName:  "Chen",
	}
	patchBody, _ := json.Marshal(userUpdateParams)

	app := fiber.New()
	app.Patch("/users/:id", suite.userHandler.HandleUpdateUser)

	req := httptest.NewRequest("PATCH", "/users/"+userId.Hex(), bytes.NewReader(patchBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	fmt.Println("resp", resp)
	assert.Equal(suite.T(), 200, resp.StatusCode)
}

func TestUserSuiteHandler(t *testing.T) {
	suite.Run(t, new(UserSuiteHandler))
}
