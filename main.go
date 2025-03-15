package main

import (
	"context"
	"flag"
	"log"

	"github.com/ctchen1999/hotel-system/api"
	"github.com/ctchen1999/hotel-system/api/middleware"
	"github.com/ctchen1999/hotel-system/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		return c.Status(code).JSON(fiber.Map{"error": err.Error()})
	},
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

	adminApi.Get("/user", userHandler.HandleGetUsers)
	adminApi.Get("/user/:id", userHandler.HandleGetUser)
	adminApi.Post("/user", userHandler.HandlePostUser)
	adminApi.Delete("/user/:id", userHandler.HandleDeleteUser)
	adminApi.Patch("/user/:id", userHandler.HandleUpdateUser)

	adminApi.Get("/hotel", hotelHandler.HandleGetHotels)
	adminApi.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	adminApi.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	adminApi.Post("/room/:id/book", roomHandler.HandleBookRoom)
	adminApi.Get("/room/booking", roomHandler.HandleGetBookings)

	app.Listen(*listenAddr)
}
