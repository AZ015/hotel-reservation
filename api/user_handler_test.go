package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"net/http/httptest"
	"testing"
)

type testDB struct {
	db.UserStore
}

func setup(t *testing.T) *testDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	return &testDB{
		UserStore: db.NewMongoUserStore(client),
	}
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@email.com",
		FirstName: "James",
		LastName:  "Test",
		Password:  "qwe123123",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return
	}

	if len(user.ID) == 0 {
		t.Errorf("expected a user id to be set")
	}

	if len(user.EncryptedPassword) == 0 {
		t.Errorf("expected a password to be set")
	}

	if err != nil {
		t.Error(err)
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected username %s, actual %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s, actual %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s, actual %s", params.Email, user.Email)
	}
}
