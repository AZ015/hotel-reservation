package api

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hotel-reservation/types"
	"net/http/httptest"
	"testing"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
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
