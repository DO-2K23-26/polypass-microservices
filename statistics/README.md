# Statistics Service

The Statistics Service is a microservice that collects, processes, and exposes statistical metrics for the Polypass application. It follows clean architecture principles to ensure separation of concerns and maintainability.

## Architecture

The service is structured according to clean architecture principles:

### Domain Layer

The core business logic and entities:

- **Models**: Defines the core entities like `Event` and `Metric`, as well as interfaces like `MetricCalculator`.
- **Repositories**: Defines interfaces for data access, such as `EventRepository` and `MetricRepository`.

### Application Layer

Contains the application business rules and use cases:

- **Services**: Implements business logic for events and metrics, including `EventService` and `MetricService`.

### Infrastructure Layer

Implements the interfaces defined in the domain layer:

- **Repositories**: Concrete implementations of repository interfaces, such as `EventStoreDBEventRepository` and `EventStoreDBMetricRepository`.
- **API**: HTTP handlers for exposing metrics and events via REST API.
- **Kafka**: Kafka consumer for ingesting events from other services.

## Features

- Event collection and storage
- Metric calculation and storage
  - User login count
  - Credential count (number of credentials a user has)
  - Credential access count (number of times credentials have been accessed)
- REST API for accessing metrics and events
- Kafka integration for event ingestion

## Extending the Service

### Adding New Event Types

1. Add the new event type to the `EventType` enum in `domain/models/event.go`.
2. Ensure the Kafka consumer is subscribed to the appropriate topic.

### Adding New Metrics

1. Implement the `MetricCalculator` interface in a new file under `domain/models/`.
2. Register the new calculator in the `main.go` file.

Example:

```go
// UserLoginCountCalculator calculates the number of user logins
type UserLoginCountCalculator struct{}

// Calculate computes the metric based on events
func (c *UserLoginCountCalculator) Calculate(events []Event) (Metric, error) {
    count := 0
    for _, event := range events {
        if event.Type == string(UserLoginEvent) {
            count++
        }
    }

    return Metric{
        ID:          uuid.New().String(),
        Name:        "user_login_count",
        Description: "Number of user logins",
        Value:       count,
        Unit:        "count",
        Timestamp:   time.Now(),
        Tags:        []string{string(UserMetric)},
    }, nil
}

// GetMetricType returns the type of metric
func (c *UserLoginCountCalculator) GetMetricType() MetricType {
    return CountMetric
}

// GetMetricName returns the name of the metric
func (c *UserLoginCountCalculator) GetMetricName() string {
    return "user_login_count"
}
```

## Configuration

The service can be configured using environment variables:

- `EVENTSTOREDB_CONNECTION_STRING`: EventStoreDB connection string (default: `esdb://localhost:2113?tls=false`)
- `KAFKA_BOOTSTRAP`: Kafka bootstrap servers (default: `localhost:9092`)
- `KAFKA_GROUP_ID`: Kafka consumer group ID (default: `statistics-service`)
- `HTTP_LISTEN_ADDRESS`: HTTP server listen address (default: `:8080`)

## API Endpoints

### Metrics

- `GET /metrics`: Get the latest metrics
- `GET /metrics/{id}`: Get a metric by ID
- `GET /metrics/name/{name}`: Get metrics by name
- `GET /metrics/category/{category}`: Get metrics by category
- `POST /metrics/calculate`: Calculate metrics for a specific time range

## Running the Service

```bash
# Build the service
go build -o statistics-service

# Run the service
./statistics-service
```

## Dependencies

- EventStoreDB for data storage
- Kafka for event ingestion
- Gorilla Mux for HTTP routing
