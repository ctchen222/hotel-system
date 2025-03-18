package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/ctchen1999/hotel-system/internal/api"
	"github.com/ctchen1999/hotel-system/internal/api/middleware"
	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if apiError, ok := err.(response.Error); ok {
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := response.NewError(http.StatusInternalServerError, err.Error())
		return c.Status(apiError.Code).JSON(apiError)
	},
	BodyLimit: 10 * 1024 * 1024,
}

func main() {
	listenAddr := flag.String("listen", ":8080", "server listen address")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// Handler initialization
	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler  = api.NewUserHandler(store)
		authHandler  = api.NewAuthHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		roomHandler  = api.NewRoomHandler(store)

		app      = fiber.New(config)
		api      = app.Group("/api")
		adminApi = app.Group("/admin/api", middleware.JWTAuthentication(userStore))
	)

	api.Post("/login", authHandler.HandleLogin)
	api.Post("/register", userHandler.HandlePostUser)

	adminApi.Get("/user", userHandler.HandleGetUsers)
	adminApi.Get("/user/:id", userHandler.HandleGetUser)
	adminApi.Delete("/user/:id", userHandler.HandleDeleteUser)
	adminApi.Patch("/user/:id", userHandler.HandleUpdateUser)

	adminApi.Post("/hotel", hotelHandler.HandlePostHotel)
	adminApi.Get("/hotel", hotelHandler.HandleGetHotels)
	adminApi.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	adminApi.Put("/hotel/:id", hotelHandler.HandleUpdateHotel)
	adminApi.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	adminApi.Post("/room/:id/book", roomHandler.HandleBookRoom)
	adminApi.Get("/room/booking", roomHandler.HandleGetBookings)

	app.Listen(*listenAddr)
}
