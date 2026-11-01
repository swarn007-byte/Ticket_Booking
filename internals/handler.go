package booking

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// ConfirmRequest carries the parameters needed to confirm or release a held
// booking. It is shared by both endpoints to keep the wire format consistent.
type ConfirmRequest struct {
	UserID string `json:"user_id"`
}

type Handler struct {
	store BookingStore
}

func NewHandler(store BookingStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /bookings", h.handleBookSeat)
	mux.HandleFunc("GET /movies/{movieID}/bookings", h.handleListBookings)
	mux.HandleFunc("POST /bookings/{bookingID}/confirm", h.handleConfirm)
	mux.HandleFunc("POST /bookings/{bookingID}/release", h.handleRelease)
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

func (h *Handler) handleConfirm(w http.ResponseWriter, r *http.Request) {
	bookingID := r.PathValue("bookingID")
	if bookingID == "" {
		http.Error(w, "Missing bookingID", http.StatusBadRequest)
		return
	}

	var req ConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := h.store.Confirm(ctx, bookingID, req.UserID)
	if err != nil {
		switch err {
		case ErrSessionsNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case ErrUnauthorized:
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleRelease(w http.ResponseWriter, r *http.Request) {
	bookingID := r.PathValue("bookingID")
	if bookingID == "" {
		http.Error(w, "Missing bookingID", http.StatusBadRequest)
		return
	}

	var req ConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.store.Release(ctx, bookingID, req.UserID); err != nil {
		switch err {
		case ErrSessionsNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case ErrUnauthorized:
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}