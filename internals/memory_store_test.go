package booking

import (
	"sync"
	"sync/atomic"
	"testing"
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
