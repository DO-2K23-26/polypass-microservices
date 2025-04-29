package credential

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
	"github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
)

type ICredentialRepository interface {
	CreateCredential(query CreateCredentialQuery) (*CreateCredentialResult, error)
	UpdateCredential(query UpdateCredentialQuery) (*UpdateCredentialResult, error)
	DeleteCredential(query DeleteCredentialQuery) error
	GetCredential(query GetCredentialQuery) (*GetCredentialResult, error)
	SearchCredentials(query SearchCredentialQuery) (*SearchCredentialResult, error)
	AddTagsToCredential(query AddTagsToCredentialQuery) error
	RemoveTagsFromCredential(query RemoveTagsFromCredentialQuery) error
}

type CredentialRepository struct {
	esClient infrastructure.ElasticAdapter
}

func NewCredentialRepository(esClient infrastructure.ElasticAdapter) *CredentialRepository {
	return &CredentialRepository{
		esClient: esClient,
	}
}

// Function to write a credential to elasticsearch
func (c *CredentialRepository) CreateCredential(query CreateCredentialQuery) (*CreateCredentialResult, error) {

	_, err := c.esClient.Client.Index(types.CredentialIndex).Id(query.ID).Request(query).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error indexing document ID=%s: %w", query.ID, err)
	}

	return &CreateCredentialResult{
		Credential: types.Credential{
			ID:       query.ID,
			Title:    query.Title,
			Tags:     query.Tags,
			Folder:   query.Folder,
		},
	}, nil
}

// Function to get a credential from elasticsearch
func (c *CredentialRepository) GetCredential(query GetCredentialQuery) (*GetCredentialResult, error) {
	if query.ID == "" {
		return nil, fmt.Errorf("ID is required")
	}

	res, err := c.esClient.Client.Get(types.CredentialIndex, query.ID).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting document ID=%s: %w", query.ID, err)
	}

	if !res.Found { // Document not found
		return nil, fmt.Errorf("document ID=%s not found", query.ID)
	}

	var credential types.Credential
	if err := json.Unmarshal(res.Source_, &credential); err != nil {
		return nil, fmt.Errorf("error unmarshalling document ID=%s: %w", query.ID, err)
	}

	return &GetCredentialResult{
		Credential: credential,
	}, nil
}
