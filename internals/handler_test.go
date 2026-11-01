package booking

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubStore struct {
	bookings []Booking
}

func (s *stubStore) Book(b Booking) (Booking, error) {
	return b, nil
}

func (s *stubStore) ListBookings(movieID string) []Booking {
	var out []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieID {
			out = append(out, b)
		}
	}
	return out
}

func (s *stubStore) Confirm(ctx context.Context, sessionID, userID string) (Booking, error) {
	if sessionID == "missing" {
		return Booking{}, ErrSessionsNotFound
	}
	return Booking{ID: sessionID, Status: StatusConfirmed}, nil
}

func (s *stubStore) Release(ctx context.Context, sessionID, userID string) error {
	return nil
}

func TestHandleBookSeat_MissingFields(t *testing.T) {
	h := NewHandler(&stubStore{})

	req := httptest.NewRequest(http.MethodPost, "/bookings", strings.NewReader(`{"movie_id":"m1"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.handleBookSeat(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleBookSeat_Success(t *testing.T) {
	h := NewHandler(&stubStore{})

	body := `{"movie_id":"m1","seat_id":"A1","user_id":"u1"}`
	req := httptest.NewRequest(http.MethodPost, "/bookings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.handleBookSeat(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var b Booking
	if err := json.NewDecoder(w.Body).Decode(&b); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if b.MovieID != "m1" || b.SeatID != "A1" {
		t.Fatalf("unexpected booking: %+v", b)
	}
}

func TestHandleListBookings(t *testing.T) {
	h := NewHandler(&stubStore{bookings: []Booking{
		{MovieID: "m1", SeatID: "A1"},
		{MovieID: "m1", SeatID: "A2"},
		{MovieID: "m2", SeatID: "B1"},
	}})

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	req := httptest.NewRequest(http.MethodGet, "/movies/m1/bookings", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var bs []Booking
	if err := json.NewDecoder(w.Body).Decode(&bs); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(bs) != 2 {
		t.Fatalf("expected 2 bookings for m1, got %d", len(bs))
	}
}

func TestHandleConfirm_UnknownBooking(t *testing.T) {
	h := NewHandler(&stubStore{})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/bookings/missing/confirm", strings.NewReader(`{"user_id":"u1"}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandleConfirm_Success(t *testing.T) {
	h := NewHandler(&stubStore{})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"user_id":"u1"}`
	req := httptest.NewRequest(http.MethodPost, "/bookings/abc/confirm", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var b Booking
	if err := json.NewDecoder(w.Body).Decode(&b); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if b.Status != StatusConfirmed {
		t.Fatalf("expected StatusConfirmed, got %s", b.Status)
	}
}

func TestHandleRelease_Success(t *testing.T) {
	h := NewHandler(&stubStore{})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := `{"user_id":"u1"}`
	req := httptest.NewRequest(http.MethodPost, "/bookings/abc/release", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandleRegisterRoutes(t *testing.T) {
	h := NewHandler(&stubStore{})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	cases := []struct {
		method, path string
	}{
		{"POST", "/bookings"},
		{"GET", "/movies/m1/bookings"},
		{"POST", "/bookings/abc/confirm"},
		{"POST", "/bookings/abc/release"},
	}
	for _, c := range cases {
		req := httptest.NewRequest(c.method, c.path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code == http.StatusNotFound {
			t.Fatalf("route %s %s not registered", c.method, c.path)
		}
	}
}