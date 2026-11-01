# concurrent_ticket

This project is a concurrent ticket booking system implemented in Go following clean architecture principles.

## Architecture

- **Domain Layer**: Contains the business logic and data models.
- **Application Layer**: Handles the business logic and interacts with the domain layer.
- **Infrastructure Layer**: Contains the implementation details for external systems and dependencies.

## Concurrency Model

The system uses `sync.RWMutex` for thread-safe operations on the memory store.

## API Endpoints

- **POST /bookings**: Books a seat.
- **GET /movies/{movieID}/bookings**: Lists all bookings for a specific movie.

## Running Tests

To run the tests with race detection, use the following command:

