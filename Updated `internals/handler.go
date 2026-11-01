package booking

import (
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	store BookingStore
}

func NewHandler(store BookingStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /bookings", h.handleBookSeat)
	mux.HandleFunc("GET /movies/{movieID}/bookings", h.handleListBookings)
}

type BookRequest struct {
	MovieID string `json:"movie_id"`
	SeatID  string `json:"seat_id"`
	UserID  string `json:"user_id"`
}

func (h *Handler) handleBookSeat(w http.ResponseWriter, r *http.Request) {
	var req BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.MovieID == "" || req.SeatID == "" || req.UserID == "" {
		http.Error(w, "Missing movie_id, seat_id, or user_id", http.StatusBadRequest)
		return
	}

	booking := Booking{
		MovieID: req.MovieID,
		SeatID:  req.SeatID,
		UserID:  req.UserID,
		Status:  StatusHeld,
	}

	result, err := h.store.Book(booking)
	if err != nil {
		if err == SeatAlreadyAllocated {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleListBookings(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	if movieID == "" {
		http.Error(w, "Missing movieID", http.StatusBadRequest)
		return
	}

	bookings := h.store.ListBookings(movieID)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
