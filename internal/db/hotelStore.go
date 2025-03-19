package db

import (
	"context"
	"fmt"

	"github.com/ctchen1999/hotel-system/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	Create(context.Context, *types.Hotel) (*types.Hotel, error)
	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(ctx context.Context, params types.HotelUpdateParams, id string) error
	GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error)
	GetHotelById(context.Context, string) (*types.HotelEmbed, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) Create(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	var existedHotel types.Hotel
	err := s.coll.FindOne(ctx, bson.M{"name": hotel.Name}).Decode(&existedHotel)
	if err == nil {
		return nil, fmt.Errorf("hotel already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.Id = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.Id = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, params types.HotelUpdateParams, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"name":     params.Name,
			"location": params.Location,
			"rating":   params.Rating,
		},
	}

	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return err
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) GetHotelById(ctx context.Context, id string) (*types.HotelEmbed, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// var hotel []*types.Hotel
	var hotel *types.HotelEmbed
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"_id": oid}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "rooms"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "hotelId"},
			{Key: "as", Value: "rooms"},
		}}},
		bson.D{{Key: "$match", Value: bson.M{"rooms": bson.M{"$ne": bson.A{}}}}},
	}
	cur, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if cur.Next(ctx) {
		if err := cur.Decode(&hotel); err != nil {
			return nil, err
		}
	}

	return hotel, nil
}
