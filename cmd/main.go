package main

import (
	"flag"
	"net/http"

	"github.com/ctchen1999/hotel-system/internal/api"
	"github.com/ctchen1999/hotel-system/internal/api/middleware"
	"github.com/ctchen1999/hotel-system/internal/db"
	models "github.com/ctchen1999/hotel-system/internal/pg"
	"github.com/ctchen1999/hotel-system/internal/response"
	"github.com/gofiber/fiber/v2"
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

	client := db.NewMongoInstance(db.MONGOURI)
	defer client.Disconnect(db.Ctx)
	pool := models.NewPostgresInstance(models.Ctx, models.PGURI)
	pool.DB.Ping(db.Ctx)
	defer pool.DB.Close()

	// Handler initialization
	var (
		pgUserStore  = models.NewPostgresUserStore(pool)
		pgHotelStore = models.NewPostgresHotelStore(pool)

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

		pgUserHandler  = api.NewPgUserHandler(pgUserStore)
		pgHotelHandler = api.NewPgHotelHandler(pgHotelStore)
		pgAuthHandler  = api.NewPgAuthHandler(pgUserStore)

		userHandler  = api.NewUserHandler(store)
		authHandler  = api.NewAuthHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		roomHandler  = api.NewRoomHandler(store)

		app        = fiber.New(config)
		api        = app.Group("/api")
		adminPgApi = app.Group("/admin/pg", middleware.PgJWTAuthentication(pgUserStore))
		adminApi   = app.Group("/admin/api", middleware.MongoJWTAuthentication(userStore))
	)

	// MONGODB
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

	// POSTGRES
	api.Post("/pg/login", pgAuthHandler.HandleLogin)

	adminPgApi.Get("/user", pgUserHandler.HandleGetUsers)
	adminPgApi.Get("/user/:id", pgUserHandler.HandleGetUser)
	adminPgApi.Delete("/user/:id", pgUserHandler.HandleDeleteUser)
	adminPgApi.Post("/user", pgUserHandler.HandleCreateUser)
	adminPgApi.Patch("/user/:id", pgUserHandler.HandleUpdateUser)

	adminPgApi.Post("/hotel", pgHotelHandler.HandleCreateHotel)
	adminPgApi.Get("/hotel", pgHotelHandler.HandleGetHotels)
	adminPgApi.Get("/hotel/:id", pgHotelHandler.HandleGetHotel)
	adminPgApi.Patch("/hotel/:id", pgHotelHandler.HandleUpdateHotel)
	adminPgApi.Delete("/hotel/:id", pgHotelHandler.HandlerDeleteHotel)

	app.Listen(*listenAddr)
}
