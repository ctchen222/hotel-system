package db

import (
	"context"
	"fmt"

	"github.com/ctchen1999/hotel-system/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper

	GetUserById(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	Create(context.Context, *types.User) (*types.User, error)
	DeleteById(context.Context, string) error
	Update(ctx context.Context, params types.UserUpdateParams, id string) error
}

type PostgresUserStore struct{}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	// validate the correctness of the id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// find the user by id
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	// find the user by email
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) Create(ctx context.Context, user *types.User) (*types.User, error) {
	// First check if a user with the same email exists
	var existedUser types.User
	err := s.coll.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existedUser)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", existedUser.Email)
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// If no user with same email exists, proceed with insertion
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Id = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) DeleteById(ctx context.Context, id string) error {
	// validate the correctness of the id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) Update(ctx context.Context, params types.UserUpdateParams, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": oid,
	}

	update := bson.M{
		"$set": bson.M{
			"firstName": params.FirstName,
			"lastName":  params.LastName,
		},
	}

	if _, err := s.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("---Dropping User collection---")
	return s.coll.Drop(ctx)
}
