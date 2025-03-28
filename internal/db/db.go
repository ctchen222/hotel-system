package db

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBNAME      = "hotel-reservation"
	DBTESTNAME  = "hotel-reservation-test"
	MONGOURI    = "mongodb://localhost:27017"
	userColl    = "users"
	hotelColl   = "hotels"
	roomColl    = "rooms"
	bookingColl = "bookings"
)

var (
	MongoInstance *mongo.Client
	mongoOnce     sync.Once
	Ctx           = context.Background()
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func ToObjectId(id string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return objectId
}

func NewMongoInstance(connString string) *mongo.Client {
	mongoOnce.Do(func() {
		client, err := mongo.Connect(Ctx, options.Client().ApplyURI(connString))
		if err != nil {
			log.Fatal(err)
		}
		MongoInstance = client
	})
	return MongoInstance
}
