package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"hotel-reservation/db/fixtures"
	"log"
	"math/rand"
	"time"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   db.NewMongoHotelStore(client),
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Printf("%s -> %s\n", user.FirstName, api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Printf("%s -> %s\n", admin.FirstName, api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(store, "medium", true, 45.67, hotel.ID)
	fmt.Println(room)

	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 3))
	fmt.Println(booking)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random_hotel %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
