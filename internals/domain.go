package booking

import (
	"context"
	"errors"
	"time"
)

var SeatAlreadyAllocated = errors.new("seat is already allocated")
var ErrSessionsNotFound = errors.new("sessions not found")


type Booking struct{
	ID string
	MovieID string
	SeatID string
	UserID string
	Status string
	ExpiresAt string
}

// type booking struct{
// 	id string,
// 	movie_id string,
// 	seat_id string,
// 	user_id string,
// 	status string
// }
type BookingStore interface {
	Book(b Booking) (Booking, error)                                        
	ListBookings(movieID string) []Booking                                  
	Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) 
	Release(ctx context.Context, sessionID string, userID string) error    
}

// type booking_store interface{
// 	book(b booking)error
// 	list_bookings(movie_id string)[]booking
// 	list_all_detail(id string)[]booking
// }