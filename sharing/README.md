## Launch the microservice

1. ### Start MongoDB

```bash
docker-compose up -d
```

2. ### Start the Go application

```bash
go run main.go
```

## Environment Variables

This microservice requires the following environment variables to be set in a `.env` file:

- `MONGO_URI`: MongoDB connection string
- `MONGO_DB_NAME`: Name of the MongoDB database
- `PORT`: Port on which the service will run

You can create a `.env` file in the root directory of the project with these variables before launching the service.
