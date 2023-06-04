package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"time"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	bookingStore db.BookingStore
	userStore    db.UserStore
	ctx          = context.Background()
)

func seedUser(fName, lName, email, password string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fName,
		LastName:  lName,
		Password:  password,
	})

	user.IsAdmin = isAdmin

	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))

	return insertedUser
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Price:   price,
		HotelID: hotelID,
		Seaside: ss,
	}

	insertedRoom, err := roomStore.InsertRoom(ctx, room)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("room -> %s\n", insertedRoom.ID)

	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}

	insertedBooking, err := bookingStore.InsertBooking(ctx, booking)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("booking -> %s\n", insertedBooking.ID)
}

func main() {
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	james := seedUser("james", "foo", "james@oo.com", "supersecurepassword", false)
	seedUser("admin", "admin", "admin@admin.com", "admin", true)

	h1 := seedHotel("Belucia", "France", 3)
	h2 := seedHotel("The cozy hotel", "The Nederland", 4)
	h3 := seedHotel("Dont die in your sleep", "London", 1)
	srh1 := seedRoom("small", true, 56.5, h1.ID)
	seedRoom("medium", true, 89.45, h1.ID)
	seedRoom("large", false, 120.00, h1.ID)

	seedRoom("medium", false, 92.45, h2.ID)
	seedRoom("small", true, 42.65, h3.ID)

	seedBooking(james.ID, srh1.ID, time.Now(), time.Now().AddDate(0, 0, 3))
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	userStore = db.NewMongoUserStore(client)
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	bookingStore = db.NewMongoBookingStore(client)
}
