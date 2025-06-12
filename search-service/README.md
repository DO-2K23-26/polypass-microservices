# Search service

## Build the CLI

Go in the `search-service` directory

```sh
cd search-service
```

```bash
go build -o search-cli cli.go
```

## Run the CLI

Making search queries:

```bash
./search-cli search-folders --query "example" --limit 10 --page 1 --user-id "user123"
```

```bash
./search-cli search-tags --query "example" --folder-id "folder123" --limit 10 --page 1 --user-id "user123"
```

```bash
./search-cli search-credentials --query "example" --folder-id "folder123" --tag-ids "tag1,tag2" --limit 10 --page 1 --user-id "user123"
```

Producing data for the Search service (via Kafka):

```bash
KAFKA_SERVER=localhost:9092 ./search-cli produce-kafka --topic tag_creation --schema tag_creation --message '{"id":"tag1","name":"Tag Name","folder_id":"folder123"}'
```

## Launch the application

## Host setup

By default, the stack exposes the following ports:

* 5044: Logstash Beats input
* 50000: Logstash TCP input
* 9600: Logstash monitoring API
* 9200: Elasticsearch HTTP
* 9300: Elasticsearch TCP transport
* 5601: Kibana

## Bringing up the stack

Go in the `docker-elk` directory

```sh
cd docker-elk
```

The first time you launch the stack, you will need to run the setup script:

```sh
docker compose up setup
```

If everything went well and the setup completed without error, start the other stack components:

```sh
docker compose up -d
```

Launch the search service:

```sh
cd ..
go run main.go
```

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

Note: Topics name could change but those described are good enough to be understandable for everyone.

## Utils

### Generate protos file

Export the go path:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Generate the proto files:

```sh
protoc --go_out=./gen  --go-grpc_out=./gen ./proto/search.proto
```

## Optique Suggestions

### Prerequisites

Optional:

- [Swaggo](https://github.com/swaggo/swag) : for generating swagger documentation
- [Migrate](https://github.com/golang-migrate/migrate) : for database migrations
- [Just](https://github.com/casey/just) : for running the application and other helpers
- [Air](https://github.com/cosmtrek/air) : for hot reloading
