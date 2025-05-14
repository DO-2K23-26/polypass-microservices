# Organization service
This service is responsible for managing folders and tags in Polypass application. It provides an API for creating, updating, and deleting folders and tags, as well as for retrieving a list of all folders and tags.

## Setup
1. Clone the repository:
```bash
git clone git@github.com:DO-2K23-26/polypass-microservices.git
cd polypass-microservices
```

2. Start the Docker containers:
```bash
docker compose up -d
```

3. Start the service:
```bash
cd organization
go run main.go
```