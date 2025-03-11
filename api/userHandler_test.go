package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/ctchen1999/hotel-system/db"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestUserStore struct {
	db.UserStore
}

func (tbc *TestUserStore) tearDown(ctx context.Context) {
	if err := tbc.Drop(ctx); err != nil {
		log.Fatal(err)
	}
}

func setup(t *testing.T) *TestUserStore {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		t.Fatal(err)
	}

	return &TestUserStore{
		db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	testUserStore := setup(t)
	// defer testUserStore.tearDown(context.Background())

	app := fiber.New()
	userHandler := NewUserHandler(testUserStore)
	app.Post("/user", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Password:  faker.Password(),
	}
	// var params types.CreateUserParams
	// err := faker.FakeData(&params)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	fmt.Printf("params fake data: %+v\n", params)

	b, _ := json.Marshal(params)

	// fmt.Println(b)
	// fmt.Println(string(b))
	req := httptest.NewRequest("POST", "/user", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if user.FirstName != params.FirstName {
		t.Fatalf("expected first name %s, got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Fatalf("expected last name %s, got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Fatalf("expected email %s, got %s", params.Email, user.Email)
	}
}
