package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	DBNAME     = "hotel-reservation"
	DBTESTNAME = "hotel-reservation-test"
	DBURI      = "mongodb://localhost:27017"
	userColl   = "users"
	hotelColl  = "hotels"
	roomColl   = "rooms"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}

func ToObjectId(id string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return objectId
}
