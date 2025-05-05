package credential

import (
	"context"
	"encoding/json"
	"fmt"

	types "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	esTypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
)

type ICredentialRepository interface {
	Create(query CreateCredentialQuery) (*CreateCredentialResult, error)
	Update(query UpdateCredentialQuery) error
	Delete(query DeleteCredentialQuery) error
	Get(query GetCredentialQuery) (*GetCredentialResult, error)
	Search(query SearchCredentialQuery) (*SearchCredentialResult, error)
	AddTags(query AddTagsToCredentialQuery) error
	RemoveTags(query RemoveTagsFromCredentialQuery) error
}

type CredentialRepository struct {
	esClient infrastructure.ElasticAdapter
}

func NewCredentialRepository(esClient infrastructure.ElasticAdapter) ICredentialRepository {
	return &CredentialRepository{
		esClient: esClient,
	}
}

// Function to write a credential to elasticsearch
func (c CredentialRepository) Create(query CreateCredentialQuery) (*CreateCredentialResult, error) {
	credential := types.Credential{
		ID:     query.ID,
		Title:  query.Title,
		Tags:   query.Tags,
		Folder: query.Folder,
	}
	err := c.esClient.CreateDocument(types.CredentialIndex, credential.ID, credential)
	if err != nil {
		return nil, err
	}
	return &CreateCredentialResult{
		Credential: credential,
	}, nil
}

// Function to get a credential by ID from elasticsearch
func (c CredentialRepository) Get(query GetCredentialQuery) (*GetCredentialResult, error) {
	if query.ID == "" {
		return nil, fmt.Errorf("ID is required")
	}
	var credential types.Credential

	err := c.esClient.GetDocument(types.CredentialIndex, query.ID, &credential)
	if err != nil {
		return nil, fmt.Errorf("Error getting folder: %w", err)
	}

	return &GetCredentialResult{
		Credential: credential,
	}, nil
}

// Function to search for credentials on title, filtered by folder and/or tags, in elasticsearch, with a paginated result.
func (c CredentialRepository) Search(query SearchCredentialQuery) (*SearchCredentialResult, error) {
	// Default limit and offset if not provided
	limit := 10
	if query.Limit != nil {
		limit = *query.Limit
	}
	offset := 0
	if query.Offset != nil {
		offset = *query.Offset
	}

	searchOnFields := []string{"title^2", "tags.name"}
	filters := []esTypes.Query{}

	if query.FoldersScope == nil || (query.FoldersScope != nil && len(*query.FoldersScope) == 0) {
		searchOnFields = append(searchOnFields, "folder.name")
	} else {
		filters = append(filters, esTypes.Query{
			Terms: &esTypes.TermsQuery{ // Filter on folder scopes.
				TermsQuery: map[string]esTypes.TermsQueryField{
					"folder.id": *query.FoldersScope,
				},
			},
		})
	}
	// Construct the search query
	searchQuery := c.esClient.Client.Search().Index(types.CredentialIndex).Request(&search.Request{
		Query: &esTypes.Query{
			Bool: &esTypes.BoolQuery{
				Must: []esTypes.Query{
					{
						MultiMatch: &esTypes.MultiMatchQuery{
							Query:  query.SearchQuery,
							Fields: searchOnFields,
							Type:   &textquerytype.Phraseprefix, // To match on parts of words (instead of whole words).
						},
					},
				},
				Filter: filters,
			},
		},
		From: &offset,
		Size: &limit,
	})

	// Execute the search query
	res, err := searchQuery.Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %w", err)
	}

	// Parse the search results
	credentials := make([]types.Credential, len(res.Hits.Hits))
	for i, hit := range res.Hits.Hits {
		if err := json.Unmarshal(hit.Source_, &credentials[i]); err != nil {
			return nil, fmt.Errorf("error unmarshalling hit source: %w", err)
		}
	}

	return &SearchCredentialResult{
		Credentials: credentials,
		Total:       int(res.Hits.Total.Value),
		Limit:       limit,
		Offset:      offset,
	}, nil
}

func (c CredentialRepository) AddTags(query AddTagsToCredentialQuery) error {
	return nil
}

func (c CredentialRepository) Delete(query DeleteCredentialQuery) error {
	return c.esClient.DeleteDocument(types.CredentialIndex, query.ID)
}

func (c CredentialRepository) RemoveTags(query RemoveTagsFromCredentialQuery) error {
	return nil
}

// UpdateCredential implements ICredentialRepository.
func (c CredentialRepository) Update(query UpdateCredentialQuery) error {
	updateCredential := types.Credential{
		Title:  *query.Title,
		Folder: query.Folder,
	}
	err := c.esClient.UpdateDocument(types.CredentialIndex, query.ID, updateCredential)
	if err != nil {
		return err
	}
	return nil
}
