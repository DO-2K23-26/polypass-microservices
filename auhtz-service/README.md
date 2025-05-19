# Authorization Service

This service is responsible for managing authorization logic using SpiceDB and Kafka consumers. It handles events related to folders, credentials, tags, and users.

## Services

The following services are implemented to handle the logic for each event type. Each service uses the `SpiceDBAdapter` to interact with SpiceDB.

### Folder Service

**Path:** `services/folder/folder_service.go`

Handles folder-related events such as:
- Creating a folder
- Deleting a folder
- Updating a folder (e.g., changing its parent)

### Credential Service

**Path:** `services/credential/credential_service.go`

Handles credential-related events such as:
- Creating a credential
- Updating a credential
- Deleting a credential

### Tag Service

**Path:** `services/tag/tag_service.go`

Handles tag-related events such as:
- Creating a tag
- Updating a tag
- Deleting a tag

### User Service

**Path:** `services/user/user_service.go`

Handles user-related events such as:
- Adding a user to a folder
- Removing a user from a folder

## Event Controllers

The event controllers in the `controllers/events` directory delegate the event handling logic to the corresponding services.

### Folder Event Controller

**Path:** `controllers/events/folder.go`

Delegates folder-related events to the `FolderService`.

### Credential Event Controller

**Path:** `controllers/events/credential.go`

Delegates credential-related events to the `CredentialService`.

### Tag Event Controller

**Path:** `controllers/events/tag.go`

Delegates tag-related events to the `TagService`.

### User Event Controller

**Path:** `controllers/events/user.go`

Delegates user-related events to the `UserService`.

## Kafka Consumers

The Kafka consumers in `internal/consumers/consumers.go` are responsible for consuming events from Kafka topics and passing them to the appropriate event controllers.

### Topics

- **Folder Events**: `create_folder`, `delete_folder`, `update_folder`
- **Credential Events**: `create_credential`, `update_credential`, `delete_credential`
- **Tag Events**: `create_tag`, `update_tag`, `delete_tag`
- **User Events**: `add_user`, `revoke_user`

## Dependencies

- **SpiceDB**: Used for managing relationships and permissions.
- **Kafka**: Used for event-driven communication.

## How to Add a New Event Type

1. Create a new service in the `services` directory.
2. Implement the logic for the new event type in the service.
3. Add a new event controller in the `controllers/events` directory.
4. Register the new Kafka topic in `internal/consumers/consumers.go`.

## Testing

- Unit tests should be written for each service to ensure the logic is correct.
- Integration tests should verify the interaction between Kafka, the event controllers, and the services.