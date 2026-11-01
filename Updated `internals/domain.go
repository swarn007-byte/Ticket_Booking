package booking

import (
	"context"
	"errors"
	"time"
)

var (
	StatusHeld      = "HELD"
	StatusBooked    = "BOOKED"
	StatusConfirmed = "CONFIRMED"
	StatusReleased  = "RELEASED"
)

var SeatAlreadyAllocated = errors.New("seat is already allocated")
var ErrSessionsNotFound = errors.New("sessions not found")

type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
	Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
	Release(ctx context.Context, sessionID string, userID string) error
}
