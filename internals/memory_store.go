package booking

import (
	"context"
	"sync"
)

type MemoryStore struct {
	mu       sync.RWMutex
	bookings map[string]Booking
}

// NewMemoryStore initialises the thread-safe memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: make(map[string]Booking),
	}
}

// Book handles thread-safe booking with a composite key (MovieID + SeatID).
func (s *MemoryStore) Book(b Booking) (Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := b.MovieID + ":" + b.SeatID
	if _, exists := s.bookings[key]; exists {
		return Booking{}, SeatAlreadyAllocated
	}

	s.bookings[key] = b
	return b, nil
}

// ListBookings returns all bookings for a specific movie thread-safely.
func (s *MemoryStore) ListBookings(movieID string) []Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieID {
			result = append(result, b)
		}
	}
	return result
}

// Confirm handles session confirmation (Milestone 4 placeholder).
func (s *MemoryStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	return Booking{}, ErrSessionsNotFound
}

// Release handles releasing a hold (Milestone 4 placeholder).
func (s *MemoryStore) Release(ctx context.Context, sessionID string, userID string) error {
	return ErrSessionsNotFound
}
