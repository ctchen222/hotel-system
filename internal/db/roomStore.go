package db

import (
	"context"

	"github.com/ctchen222/hotel-system/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	Insert(context.Context, *types.Room) (*types.Room, error)
	GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(roomColl),

		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) Insert(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.Id = resp.InsertedID.(primitive.ObjectID)

	// filter := bson.M{"_id": room.HotelId}
	// update := bson.M{"$push": bson.M{"rooms": room.Id}}
	// if err := s.HotelStore.Update(ctx, filter, update); err != nil {
	// 	return nil, err
	// }

	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var room []*types.Room
	if err := cur.All(ctx, &room); err != nil {
		return nil, err
	}
	return room, nil
}
