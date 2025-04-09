package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/ctchen222/hotel-system/internal/api"
	"github.com/ctchen222/hotel-system/internal/db"
	"github.com/ctchen222/hotel-system/internal/db/mocks"
	"github.com/ctchen222/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

type AuthSuiteHandler struct {
	suite.Suite
	mockUserStore *mocks.MockUserStore
	authHandler   *api.AuthHandler
	userHandler   *api.UserHandler

	userId primitive.ObjectID
	user   *types.User
}

func (suite *AuthSuiteHandler) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.mockUserStore = mocks.NewMockUserStore(ctrl)
	store := db.Store{
		User:    suite.mockUserStore,
		Hotel:   mocks.NewMockHotelStore(ctrl),
		Room:    mocks.NewMockRoomStore(ctrl),
		Booking: mocks.NewMockBookingStore(ctrl),
	}

	suite.authHandler = api.NewAuthHandler(suite.mockUserStore)
	suite.userHandler = api.NewUserHandler(&store)
}

func (suite *AuthSuiteHandler) BeforeTest(suiteName, testName string) {
	suite.userId = primitive.NewObjectID()
	suite.user = &types.User{
		Id:                suite.userId,
		FirstName:         "TwoBao",
		LastName:          "Chen",
		Email:             "twobao@twobao.com",
		EncryptedPassword: "$2a$10$fXSf7.i3RluVG3GMGPa7FORF2NdWB9Els7veSo13teTYXChpVHJQG",
	}
}

func (suite *AuthSuiteHandler) TestAuthHandler_HandleRegister() {
	userCreateParams := &types.CreateUserParams{
		FirstName: "TwoBao",
		LastName:  "Chen",
		Email:     "twobao@twobao.com",
		Password:  "test1234",
	}
	userBody, _ := json.Marshal(userCreateParams)

	suite.mockUserStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(suite.user, nil)

	app := fiber.New()
	app.Post("/register", suite.userHandler.HandlePostUser)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(userBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(suite.T(), 200, resp.StatusCode)
}

func (suite *AuthSuiteHandler) TestAuthHandler_HandleLogin() {
	type args struct {
		params types.AuthParams
	}
	type want struct {
		User *types.User
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "Test login success",
			args: args{
				params: types.AuthParams{
					Email:    "twobao@twbao.com",
					Password: "test1234",
				},
			},
			want: want{
				User: suite.user,
			},
			wantErr: false,
		},
	}

	suite.mockUserStore.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(suite.user, nil).Times(1)

	app := fiber.New()
	app.Post("/login", suite.authHandler.HandleLogin)
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			userBody, _ := json.Marshal(tt.args.params)

			req := httptest.NewRequest("POST", "/login", bytes.NewReader(userBody))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			if err := json.NewDecoder(resp.Body).Decode(tt.want.User); err != nil {
				fmt.Println(err)
			}
			assert.Equal(suite.T(), 200, resp.StatusCode)
			assert.Equal(suite.T(), tt.want.User, suite.user)
		})
	}
}

func TestAuthSuiteHandler(t *testing.T) {
	suite.Run(t, new(AuthSuiteHandler))
}
