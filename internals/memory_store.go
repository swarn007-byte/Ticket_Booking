package booking

import (
	"context"
	"crypto/rand"
	"sync"
	"time"
)

const idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type MemoryStore struct {
	mu       sync.RWMutex
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: make(map[string]Booking),
	}
}

func (s *MemoryStore) Book(b Booking) (Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := b.MovieID + ":" + b.SeatID
	if _, exists := s.bookings[key]; exists {
		return Booking{}, SeatAlreadyAllocated
	}

	b.ID = generateUniqueID()
	b.Status = StatusHeld
	b.ExpiresAt = time.Now().Add(5 * time.Minute)
	s.bookings[key] = b
	return b, nil
}

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

func (s *MemoryStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, b := range s.bookings {
		if b.ID != sessionID {
			continue
		}
		if b.UserID != userID {
			return Booking{}, ErrUnauthorized
		}
		if b.Status != StatusHeld {
			return Booking{}, ErrSessionsNotFound
		}
		if time.Now().After(b.ExpiresAt) {
			return Booking{}, ErrSessionsNotFound
		}

		b.Status = StatusConfirmed
		s.bookings[key] = b
		return b, nil
	}
	return Booking{}, ErrSessionsNotFound
}

func (s *MemoryStore) Release(ctx context.Context, sessionID string, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, b := range s.bookings {
		if b.ID != sessionID {
			continue
		}
		if b.UserID != userID {
			return ErrUnauthorized
		}
		if b.Status != StatusHeld {
			return ErrSessionsNotFound
		}

		delete(s.bookings, key)
		return nil
	}
	return ErrSessionsNotFound
}

func generateUniqueID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "fallback-id"
	}
	for i := range b {
		b[i] = idAlphabet[int(b[i])%len(idAlphabet)]
	}
	return string(b)
}
