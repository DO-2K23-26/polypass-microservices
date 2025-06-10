# Statistics Service Architecture

This document provides a high-level overview of the Statistics Service architecture, its components, and how they interact with each other.

## Overview

The Statistics Service is a microservice that collects, processes, and exposes statistical metrics for the Polypass application. It follows clean architecture principles to ensure separation of concerns and maintainability.

## Architecture Layers

The service is structured according to clean architecture principles, with the following layers:

### Domain Layer

The core business logic and entities:

- **Models**: Defines the core entities like `Event` and `Metric`, as well as interfaces like `MetricCalculator`.
- **Repositories**: Defines interfaces for data access, such as `EventRepository` and `MetricRepository`.

### Application Layer

Contains the application business rules and use cases:

- **Services**: Implements business logic for events and metrics, including `EventService` and `MetricService`.

### Infrastructure Layer

Implements the interfaces defined in the domain layer:

- **Repositories**: Concrete implementations of repository interfaces, such as `MongoEventRepository` and `MongoMetricRepository`.
- **API**: HTTP handlers for exposing metrics and events via REST API.
- **Kafka**: Kafka consumer for ingesting events from other services.

### Configuration

The service uses a configuration package to load configuration from a file and environment variables:

- **Config**: Defines the configuration structure and provides functions to load the configuration.

## Component Interactions

### Event Flow

1. Events are produced by other services and published to Kafka topics.
2. The Kafka consumer in the Statistics Service consumes these events.
3. The events are stored in MongoDB using the `EventRepository`.
4. The `MetricService` uses the `EventRepository` to retrieve events for metric calculation.
5. The `MetricService` uses registered `MetricCalculator` implementations to calculate metrics based on events.
6. The calculated metrics are stored in MongoDB using the `MetricRepository`.
7. The REST API exposes endpoints to retrieve metrics.

### Metric Calculation

1. The `MetricService` retrieves events from the `EventRepository`.
2. For each registered `MetricCalculator`, the `MetricService` calls the `Calculate` method with the retrieved events.
3. Each `MetricCalculator` calculates a specific metric based on the events.
4. The calculated metrics are stored in MongoDB using the `MetricRepository`.

### API Endpoints

The REST API exposes the following endpoints:

- `GET /metrics`: Get the latest metrics
- `GET /metrics/{id}`: Get a metric by ID
- `GET /metrics/name/{name}`: Get metrics by name
- `GET /metrics/category/{category}`: Get metrics by category
- `POST /metrics/calculate`: Calculate metrics for a specific time range

## Extending the Service

### Adding New Event Types

1. Add the new event type to the `EventType` enum in `domain/models/event.go`.
2. Ensure the Kafka consumer is subscribed to the appropriate topic.

### Adding New Metrics

1. Implement the `MetricCalculator` interface in a new file under `domain/models/`.
2. Register the new calculator in the `main.go` file.

## Deployment

The service is containerized using Docker and can be deployed using Docker Compose or Kubernetes. The service requires the following dependencies:

- MongoDB for data storage
- Kafka for event ingestion

## Configuration

The service can be configured using environment variables or a configuration file. The configuration includes settings for:

- MongoDB connection
- Kafka connection
- HTTP server
- Logging
- Metrics calculation