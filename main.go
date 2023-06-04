package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/api"
	"hotel-reservation/api/middleware"
	"hotel-reservation/db"
	"log"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// store init
	userStore := db.NewMongoUserStore(client)
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	bookingStore := db.NewMongoBookingStore(client)

	store := db.Store{
		Hotel:   hotelStore,
		Room:    roomStore,
		User:    userStore,
		Booking: bookingStore,
	}

	// handlers init
	authHandler := api.NewAuthHandler(userStore)
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(store)
	roomHandler := api.NewRoomHandler(store)
	bookingHandler := api.NewBookingHandler(store)

	app := fiber.New(config)
	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	admin := apiv1.Group("/admin", middleware.AdminAuth)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// rooms handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	// TODO: cancel booking

	// booking handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	// admin handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)

	err = app.Listen(*listenAddr)
	if err != nil {
		return
	}
}
