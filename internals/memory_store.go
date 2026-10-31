package booking

// type memorystore struct{
// 	bookings map[string]bookings
// }

// func newmemorystore() *memorystore{
// 	return &memorystore{
// 		bookings map[string]bookings{}
// 	}
// }

// func() book(b booking)error
// func() list_bookings(movie_id string)[]booking
// func() list_all_detail(id string)[]booking

import context


type MemoryStore struct{
	bookings map[string]Booking
}

//1st function to initialise memory store.
func NewMemoryStore *MemoryStore{
	// bookings map[string]booking
	return &MemoryStore{
		bookings: map[string]Booking{}, // empty map, ready to use
	}
}

//2nd function to book ticket.
func (s *MemoryStore) Book(b Booking) (Booking, error) {
	if _, exists := s.bookings[b.SeatID]; exists {
		return Booking{}, ErrSeatAlreadyBooked
	}
	s.bookings[b.SeatID] = b
	return b, nil
}







  