# concurrent_ticket

This project is a concurrent ticket booking system implemented in Go following clean architecture principles.

## Architecture

- **Domain Layer**: Contains the business logic and data models.
- **Application Layer**: Handles the business logic and interacts with the domain layer.
- **Infrastructure Layer**: Contains the implementation details for external systems and dependencies.

## Concurrency Model

The system uses `sync.RWMutex` for thread-safe operations on the memory store. Concurrent requests for the same seat are resolved atomically under the write lock, so exactly one request wins and the rest receive `SeatAlreadyAllocated`.

## Booking Lifecycle

1. **Hold** — `POST /bookings` reserves a seat. The booking is created with status `HELD`, a unique ID, and a 5-minute expiry.
2. **Confirm** — `POST /bookings/{bookingID}/confirm` transitions the hold to `CONFIRMED`.
3. **Release** — `POST /bookings/{bookingID}/release` cancels the hold and frees the seat.

Only the owning user may confirm or release a hold, and only while it is still `HELD` and not expired.

## API Endpoints

- **POST /bookings**: Book a seat. Body: `{ "movie_id", "seat_id", "user_id" }`.
- **GET /movies/{movieID}/bookings**: List all bookings for a specific movie.
- **POST /bookings/{bookingID}/confirm**: Confirm a held booking. Body: `{ "user_id" }`.
- **POST /bookings/{bookingID}/release**: Release a held booking. Body: `{ "user_id" }`.

## Running Tests

To run the tests with race detection, use the following command:

```bash
go test -race ./...
```

This will help ensure the system is free of race conditions.

