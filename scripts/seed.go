package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ctchen1999/hotel-system/db"
	"github.com/ctchen1999/hotel-system/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 100,
		},
		{
			Size:  "normal",
			Price: 150,
		},
		{
			Size:    "normal",
			SeaSide: true,
			Price:   200,
		},
		{
			Size:  "kingsize",
			Price: 250,
		},
	}
	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelId = insertedHotel.Id
		insertedRoom, err := roomStore.Insert(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted room", insertedRoom)
	}
}

func main() {
	fmt.Println("Seeding the database")

	seedHotel("LakeShore", "Taiwan", 5)
	seedHotel("OceanView", "Taiwan", 3)
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
