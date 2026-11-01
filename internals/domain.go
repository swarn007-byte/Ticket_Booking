package booking

import (
	"context"
	"errors"
	"time"
)

// Booking lifecycle statuses.
const (
	StatusHeld      = "HELD"
	StatusBooked    = "BOOKED"
	StatusConfirmed = "CONFIRMED"
	StatusReleased  = "RELEASED"
)

var (
	SeatAlreadyAllocated = errors.New("seat is already allocated")
	ErrSessionsNotFound  = errors.New("sessions not found")
	ErrUnauthorized      = errors.New("unauthorized")
)

// Booking is the domain entity. ExpiresAt is a time.Time so that held seats
// can be expired deterministically; Status is one of the constants above.
type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt time.Time
}

// BookingStore is the persistence contract for the booking domain.
type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
	Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
	Release(ctx context.Context, sessionID string, userID string) error
}