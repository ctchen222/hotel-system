package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	Id       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int                  `bson:"rating" json:"rating"`
}

type HotelQuery struct {
	Rooms  bool
	Rating int
}

type CreateHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}

type HotelUpdateParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}
