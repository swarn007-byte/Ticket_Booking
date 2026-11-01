package booking

import (
	"context"
	"errors"
)

var SeatAlreadyAllocated = errors.New("seat is already allocated")
var ErrSessionsNotFound = errors.New("sessions not found")

type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt string
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
	Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
	Release(ctx context.Context, sessionID string, userID string) error
}
