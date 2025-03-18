package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Size    string             `bson:"size" json:"size"`
	SeaSide bool               `bson:"seaside" json:"seaside"`
	Price   float64            `bson:"price" json:"price"`
	HotelId primitive.ObjectID `bson:"hotelId" json:"hotelId"`
}

type RoomType int

const (
	_ RoomType = iota
	SingleRoom
	DoubleRoom
	SeasideRoom
	DeluxRoom
)
