package credential

import (
	"encoding/json"
	"fmt"

	"slices"

	types "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
	esTypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
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
		Folder: &query.Folder,
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

	if query.FoldersScope == nil {
		return nil, fmt.Errorf("Error: you must provide a folder scope to search on credentials")
	}

	searchOnFields := []string{"title^2"}
	filters := []esTypes.Query{}

	// Search on the folder name
	searchOnFields = append(searchOnFields, "folder.name")
	// Filter only the folder the user has access to
	filters = append(filters, esTypes.Query{
		Terms: &esTypes.TermsQuery{ // Filter on folder scopes.
			TermsQuery: map[string]esTypes.TermsQueryField{
				"folder.id": query.FoldersScope,
			},
		},
	},
	)

	additionalQuery := esTypes.Query{}
	if query.TagIds == nil || len(*query.TagIds) == 0 {
		// This allow to search on a nested field
		// In our case we are searching on the tag name
		additionalQuery = esTypes.Query{
			Nested: &esTypes.NestedQuery{
				Path: "tags",
				Query: esTypes.Query{
					Match: map[string]esTypes.MatchQuery{
						"tags.name": {
							Query: query.SearchQuery,
						},
					},
				},
			},
		}
	} else {
		filters = append(filters, esTypes.Query{
			Nested: &esTypes.NestedQuery{
				Path: "tags",
				Query: esTypes.Query{
					Terms: &esTypes.TermsQuery{
						TermsQuery: map[string]esTypes.TermsQueryField{
							"tags.id": query.TagIds,
						},
					},
				},
			},
		},
		)
	}

	// Construct the search query
	res, total, err := c.esClient.Search(
		types.CredentialIndex,
		query.SearchQuery,
		searchOnFields,
		filters,
		additionalQuery,
	)
	if err != nil {
		return nil, err
	}
	result := make([]types.Credential, *total)
	// Parse the search results
	for i, hit := range res {
		if err := json.Unmarshal(hit, &result[i]); err != nil {
			return nil, fmt.Errorf("error unmarshalling hit source: %w", err)
		}
	}
	return &SearchCredentialResult{
		Credentials: result,
		Total:       *total,
		Limit:       limit,
		Offset:      offset,
	}, nil
}

func (c CredentialRepository) AddTags(query AddTagsToCredentialQuery) error {
	res, err := c.Get(GetCredentialQuery{ID: query.ID})
	if err != nil {
		return err
	}
	for _, tag := range query.Tags {
		res.Credential.Tags = append(res.Credential.Tags, tag)
	}
	c.esClient.UpdateDocument(types.CredentialIndex, query.ID, res.Credential)
	return nil
}

func (c CredentialRepository) Delete(query DeleteCredentialQuery) error {
	return c.esClient.DeleteDocument(types.CredentialIndex, query.ID)
}

func (c CredentialRepository) RemoveTags(query RemoveTagsFromCredentialQuery) error {
	res, err := c.Get(GetCredentialQuery{ID: query.ID})
	if err != nil {
		return err
	}

	// Filter out the tags to be removed
	updatedTags := []types.Tag{}
	for _, tag := range res.Credential.Tags {
		if !slices.Contains(query.TagIds, tag.ID) {
			updatedTags = append(updatedTags, tag)
		}
	}

	// Update the credential with the filtered tags
	res.Credential.Tags = updatedTags
	err = c.esClient.UpdateDocument(types.CredentialIndex, query.ID, res.Credential)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCredential implements ICredentialRepository.
func (c CredentialRepository) Update(query UpdateCredentialQuery) error {

	updateCredential := types.Credential{}
	if query.Title != nil {
		updateCredential.Title = *query.Title
	}
	if query.Folder != nil {
		updateCredential.Folder = query.Folder
	}
	
	err := c.esClient.UpdateDocument(types.CredentialIndex, query.ID, updateCredential)
	if err != nil {
		return err
	}
	return nil
}
