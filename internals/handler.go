package booking

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	store *MemoryStore
}

func NewHandler(store *MemoryStore) *Handler {
	return &Handler{store: store}
}

//register routes

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	//here we book a seat
	mux.HandleFunc("POST /bookings", h.handleBookSeat)

	//list bookings for a movie
	mux.HandleFunc("GET /movies/{movieID}/bookings", h.handleListBookings)
}

type BookRequest struct {
	MovieID string `json:"movie_id"`
	SeatID  string `json:"seat_id"`
	UserID  string `json:"user_id"`
}

func (h *Handler) handleBookSeat(w http.ResponseWriter, r *http.Request) {
	var req BookRequest
	// Decode incoming JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	booking := Booking{
		MovieID: req.MovieID,
		SeatID:  req.SeatID,
		UserID:  req.UserID,
		Status:  "BOOKED",
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

	// Return success JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w.Encode(result))
}
