# Search service

## Data aggregation

First of all, the search service of polypass will need to listen to events from
multiple services such as:

- Organization
- Credentials


### Organization Events

For the organization events, we will need to listen to tags events and the following topics:

- tag_creation
- tag_deletion
- tag_update

For folders events, we will need to listen to the following topics:

- folder_creation
- folder_deletion
- folder_update

For credentials events, we will need to listen to the following topics:

- credential_creation
- credential_deletion
- credential_update


We will also need to know in which folder is a user in order to not send data about folder that he does not have the access to :

- user_folder

Note: Topics name could change but those describe good enough to be understandable for everyone.

## Optique Suggestions

### Prerequisites

Optional:

- [Swaggo](https://github.com/swaggo/swag) : for generating swagger documentation
- [Migrate](https://github.com/golang-migrate/migrate) : for database migrations
- [Just](https://github.com/casey/just) : for running the application and other helpers
- [Air](https://github.com/cosmtrek/air) : for hot reloading
