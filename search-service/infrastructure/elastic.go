package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	commonTypes "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
)

type ElasticAdapter struct {
	Client *elasticsearch.TypedClient
}

func NewElasticAdapter(host string, username, password string) (*ElasticAdapter, error) {

	esConfig := elasticsearch.Config{Addresses: []string{host}, Username: username,
		Password: password,
	}

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
func (e *ElasticAdapter) CheckHealth() bool {
	_, err := e.Client.Ping().Do(context.Background())
	if err != nil {
		log.Println("Elastic health problem:", err)
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

	exists, err := e.Client.Indices.Exists(indexName).Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
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

func (e *ElasticAdapter) GetDocument(indexName string, documentId string, document any) error {
	res, err := e.Client.Get(indexName, documentId).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error getting document ID=%s: %w", documentId, err)
	}
	if !res.Found { // Document not found
		return fmt.Errorf("document ID=%s not found", documentId)
	}

	if err := json.Unmarshal(res.Source_, document); err != nil {
		return fmt.Errorf("error unmarshalling document ID=%s: %w", documentId, err)
	}
	return nil
}

func (e *ElasticAdapter) CreateDocument(indexName string, documentId string, document any) error {
	_, err := e.Client.Index(indexName).Id(documentId).Request(document).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error creating document in index %s ID=%s: %w", indexName, documentId, err)
	}
	return nil
}

func (e *ElasticAdapter) UpdateDocument(indexName string, documentId string, document any) error {
	_, err := e.Client.Update(indexName, documentId).Doc(document).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error updating document in index %s ID=%s: %w", indexName, documentId, err)
	}
	return nil
}

func (e *ElasticAdapter) DeleteDocument(indexName string, documentId string) error {
	_, err := e.Client.Delete(indexName, documentId).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting document in index %s ID=%s: %w", indexName, documentId, err)
	}
	return nil
}

func (e *ElasticAdapter) Search(
	indexName, searchQueryString string,
	searchOnField []string,
	filters []types.Query,
	additionalQueries ...types.Query,
) ([]json.RawMessage, *int, error) {
	queries := []types.Query{
		{
			MultiMatch: &types.MultiMatchQuery{
				Query:  searchQueryString,
				Fields: searchOnField,
				Type:   &textquerytype.Phraseprefix,
			},
		},
	}

	if additionalQueries == nil || len(additionalQueries) == 0 {
		additionalQueries = []types.Query{}
		for _, query := range additionalQueries {
			queries = append(queries, query)
		}

	}

	// Building the base query with the field to search on and the search Query

	searchQuery := e.Client.Search().Index(indexName).Request(&search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Should:             queries,
				MinimumShouldMatch: 1,
				Filter:             filters,
			},
		},
	})

	// Execute the search query
	res, err := searchQuery.Do(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("error executing search query: %w", err)
	}

	totalHits := int(res.Hits.Total.Value)
	result := make([]json.RawMessage, 0, totalHits)
	// Parse the search results
	hits := res.Hits.Hits
	for _, hit := range hits {
		result = append(result, hit.Source_)
	}
	return result, &totalHits, nil
}
