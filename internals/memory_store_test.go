package booking

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentBook_OnlyOneWins(t *testing.T) {
	store := NewMemoryStore()

	n := 100 // 100 goroutines trying to book the same seat at once
	var (
		winners atomic.Int64
		losers  atomic.Int64
		wg      sync.WaitGroup
	)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()

			_, err := store.Book(Booking{
				MovieID: "dune",
				SeatID:  "A1",
				UserID:  "user",
			})

			if err == nil {
				winners.Add(1)
			} else {
				losers.Add(1)
			}
		}(i)
	}

	wg.Wait()

	// Assertions
	if winners.Load() != 1 {
		t.Errorf("expected exactly 1 winner, got %d", winners.Load())
	}

	if losers.Load() != int64(n-1) {
		t.Errorf("expected %d losers, got %d", n-1, losers.Load())
	}
}

func TestBook_AssignsUniqueIDAndExpiry(t *testing.T) {
	store := NewMemoryStore()

	b1, err := store.Book(Booking{MovieID: "m1", SeatID: "A1", UserID: "u1"})
	if err != nil {
		t.Fatalf("Book failed: %v", err)
	}
	b2, err := store.Book(Booking{MovieID: "m1", SeatID: "A2", UserID: "u2"})
	if err != nil {
		t.Fatalf("Book failed: %v", err)
	}

	if b1.ID == "" {
		t.Fatal("booking ID must not be empty")
	}
	if b1.ID == b2.ID {
		t.Fatalf("booking IDs must be unique: %s", b1.ID)
	}
	if b1.Status != StatusHeld {
		t.Fatalf("expected StatusHeld, got %s", b1.Status)
	}
	if b1.ExpiresAt.Before(time.Now()) {
		t.Fatal("expiry must be in the future")
	}
	if b1.ExpiresAt.Sub(time.Now()) > 6*time.Minute {
		t.Fatal("expiry must be within 5 minutes")
	}
}

func TestConfirm_TransitionsToConfirmed(t *testing.T) {
	store := NewMemoryStore()

	hold, err := store.Book(Booking{MovieID: "m1", SeatID: "A1", UserID: "u1"})
	if err != nil {
		t.Fatalf("Book failed: %v", err)
	}

	confirmed, err := store.Confirm(context.Background(), hold.ID, "u1")
	if err != nil {
		t.Fatalf("Confirm failed: %v", err)
	}
	if confirmed.Status != StatusConfirmed {
		t.Fatalf("expected StatusConfirmed, got %s", confirmed.Status)
	}
}

func TestConfirm_UnauthorizedUser(t *testing.T) {
	store := NewMemoryStore()

	hold, err := store.Book(Booking{MovieID: "m1", SeatID: "A1", UserID: "u1"})
	if err != nil {
		t.Fatalf("Book failed: %v", err)
	}

	if _, err := store.Confirm(context.Background(), hold.ID, "u2"); err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestConfirm_AlreadyConfirmed(t *testing.T) {
	store := NewMemoryStore()

	hold, err := store.Book(Booking{MovieID: "m1", SeatID: "A1", UserID: "u1"})
	if err != nil {
		t.Fatalf("Book failed: %v", err)
	}
	if _, err := store.Confirm(context.Background(), hold.ID, "u1"); err != nil {
		t.Fatalf("first Confirm failed: %v", err)
	}
	if _, err := store.Confirm(context.Background(), hold.ID, "u1"); err != ErrSessionsNotFound {
		t.Fatalf("expected ErrSessionsNotFound for double confirm, got %v", err)
	}
}

func TestConfirm_UnknownBooking(t *testing.T) {
	store := NewMemoryStore()
	if _, err := store.Confirm(context.Background(), "nope", "u1"); err != ErrSessionsNotFound {
		t.Fatalf("expected ErrSessionsNotFound, got %v", err)
	}
}

func TestRelease_FreesSeat(t *testing.T) {
	store := NewMemoryStore()

	hold, err := store.Book(Booking{MovieID: "m1", SeatID: "A1", UserID: "u1"})
	if err != nil {
		t.Fatalf("Book failed: %v", err)
	}

	if err := store.Release(context.Background(), hold.ID, "u1"); err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	// The seat must be bookable again.
	if _, err := store.Book(Booking{MovieID: "m1", SeatID: "A1", UserID: "u2"}); err != nil {
		t.Fatalf("expected seat to be rebookable after release, got %v", err)
	}
}

func TestRelease_UnauthorizedUser(t *testing.T) {
	store := NewMemoryStore()

	hold, err := store.Book(Booking{MovieID: "m1", SeatID: "A1", UserID: "u1"})
	if err != nil {
		t.Fatalf("Book failed: %v", err)
	}

	if err := store.Release(context.Background(), hold.ID, "u2"); err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestConfirm_ExpiredHold(t *testing.T) {
	store := NewMemoryStore()

	hold, err := store.Book(Booking{MovieID: "m1", SeatID: "A1", UserID: "u1"})
	if err != nil {
		t.Fatalf("Book failed: %v", err)
	}

	// Backdate the expiry past "now".
	store.mu.Lock()
	hold.ExpiresAt = time.Now().Add(-time.Minute)
	store.bookings[hold.MovieID+":"+hold.SeatID] = hold
	store.mu.Unlock()

	if _, err := store.Confirm(context.Background(), hold.ID, "u1"); err != ErrSessionsNotFound {
		t.Fatalf("expected ErrSessionsNotFound for expired hold, got %v", err)
	}
}