package fixtures

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"time"
)

func AddUser(store *db.Store, fName, lName string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fName, lName),
		FirstName: fName,
		LastName:  lName,
		Password:  fmt.Sprintf("%s_%s", fName, lName),
	})

	user.IsAdmin = isAdmin

	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store *db.Store, name, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDs = rooms
	if rooms == nil {
		roomIDs = make([]primitive.ObjectID, 0)
	}

	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDs,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hid primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Price:   price,
		HotelID: hid,
		Seaside: ss,
	}

	insertedRoom, err := store.Room.InsertRoom(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("room -> %s\n", insertedRoom.ID)

	return insertedRoom
}

func AddBooking(store *db.Store, uid, rid primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}

	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("booking -> %s\n", insertedBooking.ID)

	return insertedBooking
}
