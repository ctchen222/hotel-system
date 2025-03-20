package db

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ctchen1999/hotel-system/internal/db/mocks"
	"github.com/ctchen1999/hotel-system/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/mock/gomock"
)

func TestNewMongoBookingStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingStore := mocks.NewMockBookingStore(ctrl)

	Id1 := primitive.NewObjectID()
	Id2 := primitive.NewObjectID()

	mockBookingStore.EXPECT().GetBookings(gomock.Any(), gomock.Any()).Return(
		[]*types.Booking{
			{
				Id:        Id1,
				UserId:    Id1,
				RoomId:    Id1,
				NumPerson: 1,
				From:      time.Now(),
				To:        time.Now().Add(time.Second * 10),
			},
			{
				Id:        Id2,
				UserId:    Id2,
				RoomId:    Id2,
				NumPerson: 2,
				From:      time.Now(),
				To:        time.Now().Add(time.Second * 20),
			},
		}, nil).Times(2)

	result, _ := mockBookingStore.GetBookings(context.Background(), bson.M{})
	fmt.Println(result)
	result, _ = mockBookingStore.GetBookings(context.Background(), bson.M{})
	fmt.Println(result)

}

func TestMongoBookingStore_InsertBookRoom(t *testing.T) {
	type fields struct {
		client *mongo.Client
		coll   *mongo.Collection
	}
	type args struct {
		ctx     context.Context
		booking *types.Booking
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Booking
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MongoBookingStore{
				client: tt.fields.client,
				coll:   tt.fields.coll,
			}
			got, err := s.InsertBookRoom(tt.args.ctx, tt.args.booking)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoBookingStore.InsertBookRoom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MongoBookingStore.InsertBookRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoBookingStore_GetBookings(t *testing.T) {
	type fields struct {
		client *mongo.Client
		coll   *mongo.Collection
	}
	type args struct {
		ctx    context.Context
		filter bson.M
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*types.Booking
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MongoBookingStore{
				client: tt.fields.client,
				coll:   tt.fields.coll,
			}
			got, err := s.GetBookings(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoBookingStore.GetBookings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MongoBookingStore.GetBookings() = %v, want %v", got, tt.want)
			}
		})
	}
}
