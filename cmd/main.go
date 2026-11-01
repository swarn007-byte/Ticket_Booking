package main

import (
	"fmt"
	"log"
	"net/http"

	"concurrent_ticket/internals"
)

func main() {
	store := booking.NewMemoryStore()

	handler := booking.NewHandler(store)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	fmt.Println("🚀 Server is running on http://localhost:8080...")
	
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}