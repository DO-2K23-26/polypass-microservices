# Search service

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

Go in the docker-elk directory
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

## Examples


### Manage nested field

As we will replicate in credential the folder object we will need to maintain the integrity of datas (given by AI prompt):

```go
// Example function to update folder name across all credentials
func UpdateFolderNameInCredentials(es *elasticsearch.Client, folderID string, newName string) error {
    // Define the query to find all credentials with the specific folder_id
    script := fmt.Sprintf(`
        if (ctx._source.folder != null && ctx._source.folder.id == '%s') {
            ctx._source.folder.name = '%s';
        }
    `, folderID, newName)
    
    // Create the update by query request
    req := esapi.UpdateByQueryRequest{
        Index:   []string{"credentials"},
        Body:    strings.NewReader(`{"query": {"term": {"folder_id": "` + folderID + `"}}, "script": {"source": "` + script + `"}}`),
        Refresh: "true",
    }
    
    res, err := req.Do(context.Background(), es)
    if err != nil {
        return err
    }
    defer res.Body.Close()
    
    if res.IsError() {
        // Handle error
        return fmt.Errorf("error updating credentials: %s", res.String())
    }
    
    return nil
}
```
### Exclude folders

This example demonstrates how to exclude folders from a search query to scope to a list of folders.

```go
func (r *elasticSearchFolderRepository) SearchFolder(query SearchFolderQuery) (*SearchFolderResult, error) {
    // Build the basic search query
    searchQuery := map[string]interface{}{
        "query": map[string]interface{}{},
        "from":  (query.Page - 1) * query.Size,
        "size":  query.Size,
    }
    
    // Define the main query part based on search term
    mainQuery := map[string]interface{}{
        "match": map[string]interface{}{
            "name": query.Term,
        },
    }
    
    // If we have allowed folder IDs, add a filter
    if len(query.AllowedFolderIds) > 0 {
        searchQuery["query"] = map[string]interface{}{
            "bool": map[string]interface{}{
                "must": mainQuery,
                "filter": map[string]interface{}{
                    "terms": map[string]interface{}{
                        "id": query.AllowedFolderIds,
                    },
                },
            },
        }
    } else {
        // Without folder filtering, just use the main query
        searchQuery["query"] = mainQuery
    }
    
    // Execute search and process results...
}
```

## Optique Suggestions

### Prerequisites

Optional:

- [Swaggo](https://github.com/swaggo/swag) : for generating swagger documentation
- [Migrate](https://github.com/golang-migrate/migrate) : for database migrations
- [Just](https://github.com/casey/just) : for running the application and other helpers
- [Air](https://github.com/cosmtrek/air) : for hot reloading
