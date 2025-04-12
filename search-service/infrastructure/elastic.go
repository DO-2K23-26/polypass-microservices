package infrastructure

import (
	"context"

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
	err := e.createIndex(commonTypes.FolderIndex, commonTypes.EsFolder)
	if err != nil {
		return err
	}
	err = e.createIndex(commonTypes.TagIndex, commonTypes.EsTag)
	if err != nil {
		return err
	}
	err = e.createIndex(commonTypes.CredentialIndex, commonTypes.EsCredential)
	if err != nil {
		return err
	}

	err = e.createIndex(commonTypes.UserIndex, commonTypes.EsUser)
	if err != nil {
		return err
	}
	
	return nil
}

func (e *ElasticAdapter) createIndex(indexName string, mapping map[string]types.Property) error {
	_, err := e.Client.Indices.Create(indexName).
		Request(&create.Request{Mappings: &types.TypeMapping{Properties: mapping}}). // Add an empty TypeMapping
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
