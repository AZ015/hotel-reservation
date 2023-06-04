package api

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
)

type BookingHandler struct {
	store db.Store
}

func NewBookingHandler(store db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// TODO: this needs to admon auth
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// TODO: this needs to user auth
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type:    "error",
			Message: "not authorized",
		})
	}

	return c.JSON(booking)
}
