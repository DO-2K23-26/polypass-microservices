package infrastructure

import (
	"context"
	"log"
	"sync"

	commonTypes "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type ElasticAdapter struct {
	Client *elasticsearch.TypedClient
}

func NewElasticAdapter(host string, password string) (*ElasticAdapter, error) {
	esConfig := elasticsearch.Config{Addresses: []string{host}, Password: password}
	Client, err := elasticsearch.NewTypedClient(esConfig)
	if err != nil {
		return nil, err
	}
	return &ElasticAdapter{
		Client,
	}, nil
}

// Ping checks the connection to the Elasticsearch server.
// False = failed connection , True = successful connection
func (e *ElasticAdapter) Ping() bool {
	_, err := e.Client.Ping().Do(context.Background())
	if err != nil {
		return false
	}
	return true
}

func (e *ElasticAdapter) CreateIndexes() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 4)

	indexes := []struct {
		name    string
		mapping map[string]types.Property
	}{
		{commonTypes.FolderIndex, commonTypes.EsFolder},
		{commonTypes.TagIndex, commonTypes.EsTag},
		{commonTypes.CredentialIndex, commonTypes.EsCredential},
		{commonTypes.UserIndex, commonTypes.EsUser},
	}

	for _, index := range indexes {
		wg.Add(1)
		go func(indexName string, mapping map[string]types.Property) {
			defer wg.Done()
			if err := e.createIndexIfNotExists(indexName, mapping); err != nil {
				errChan <- err
			}
		}(index.name, index.mapping)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (e *ElasticAdapter) createIndexIfNotExists(indexName string, mapping map[string]types.Property) error {
	res, err := e.Client.Indices.Exists(indexName).Do(context.Background())
	if err != nil {
		return err
	}
	if res.StatusCode == 404 {
		_, err = e.Client.Indices.Create(indexName).
			Request(&create.Request{Mappings: &types.TypeMapping{Properties: mapping}}).
			Do(context.Background())
		if err != nil {
			return err
		}
		log.Println("Index", indexName, "was created")
	} else {
		log.Println("Index", indexName, "already exists")
	}
	return nil
}
