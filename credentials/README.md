# Credentials Microservice

## Getting Started

### Prerequisites

- Go 1.18+
- Docker
- Docker Compose

Optional:

- [Swaggo](https://github.com/swaggo/swag) : for generating swagger documentation
- [Migrate](https://github.com/golang-migrate/migrate) : for database migrations
- [Just](https://github.com/casey/just) : for running the application and other helpers
- [Air](https://github.com/cosmtrek/air) : for hot reloading

## Running the migrations

Modify the `config.json` file to add the bootstrap flag to true.

```json
{
  "bootstrap": true
    ...
}
```

### Running the application

Turn the bootstrap flag to false and run the application.

```bash
docker compose up -d
air
```
