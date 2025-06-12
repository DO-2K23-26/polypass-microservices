# Organization service
This service is responsible for managing folders and tags in Polypass application. It provides an API for creating, updating, and deleting folders and tags, as well as for retrieving a list of all folders and tags.

## Setup
1. Clone the repository:
```bash
git clone git@github.com:DO-2K23-26/polypass-microservices.git
cd polypass-microservices/organization
```

2. Start the Docker containers:
```bash
docker compose up -d
```

3. Start the service:
```bash
go run apps/organization/cmd/organization/main.go
```

If you want to provide a host to the Credential Service, use: `CREDENTIAL_SERVICE_HOST=http://...`.
```bash
CREDENTIAL_SERVICE_HOST=http://127.0.0.1:4001 go run apps/organization/cmd/organization/main.go
```

## Folder credentials
The service exposes endpoints to manage the link between folders and credentials. Every operation forwards the request to the credential service defined by the `CREDENTIAL_SERVICE_HOST` environment variable.

- `GET /folders/{folderId}/credentials/{type}`
- `POST /folders/{folderId}/credentials/{type}`
- `PUT /folders/{folderId}/credentials/{type}/{credentialId}`
- `DELETE /folders/{folderId}/credentials/{type}`