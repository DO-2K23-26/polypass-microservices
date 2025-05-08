package tags

import (
	"encoding/json"
	"fmt"

	"github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
	esTypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type ESTagRepository struct {
	esClient infrastructure.ElasticAdapter
}

func NewESTagRepository(esClient infrastructure.ElasticAdapter) ITagRepository {
	return &ESTagRepository{
		esClient: esClient,
	}
}

// CreateTag creates a new tag in Elasticsearch
func (r *ESTagRepository) Create(query CreateTagQuery) (*CreateTagResult, error) {
	tag := types.Tag{
		ID:       query.ID,
		Name:     query.Name,
		FolderId: query.FolderID,
	}
	err := r.esClient.CreateDocument(types.TagIndex, tag.ID, tag)
	if err != nil {
		return nil, err
	}
	return &CreateTagResult{
		Tag: &tag,
	}, nil
}

// GetTag retrieves a tag by ID from Elasticsearch
func (r *ESTagRepository) Get(query GetTagQuery) (*GetTagResult, error) {
	if query.ID == "" {
		return nil, fmt.Errorf("ID is required")
	}
	var tag types.Tag

	err := r.esClient.GetDocument(types.TagIndex, query.ID, &tag)
	if err != nil {
		return nil, fmt.Errorf("Error getting tag: %w", err)
	}

	return &GetTagResult{
		Tag: &tag,
	}, nil
}

// UpdateTag updates an existing tag in Elasticsearch
// Will only update the field that are passed in the query
func (r *ESTagRepository) Update(query UpdateTagQuery) error {
	err := r.esClient.UpdateDocument(types.TagIndex, query.ID, query)
	if err != nil {
		return err
	}

	// Update the tag in all credentials
	script := `
if (ctx._source.tags != null) {
  for (tag in ctx._source.tags) {
    if (tag.id == params.tag_id) {
      if (params.containsKey('new_name') && params.new_name != null) {
        tag.name = params.new_name;
      }
      if (params.containsKey('new_folder_id') && params.new_folder_id != null) {
        tag.folder_id = params.new_folder_id;
      }
    }
  }
}`
	updateQuery := esTypes.Query{
		Nested: &esTypes.NestedQuery{
			Path: "tags",
			Query: esTypes.Query{
				Term: map[string]esTypes.TermQuery{
					"tags.id": {Value: query.ID},
				},
			},
		},
	}
	params := make(map[string]json.RawMessage)

	// Always include tag_id
	
	tagIDJson, _ := json.Marshal(query.ID)
	params["tag_id"] = tagIDJson

	// Optional fields
	if query.Name != nil {
		newNameJson, _ := json.Marshal(*query.Name)
		params["new_name"] = newNameJson
	}

	if query.FolderId != nil {
		newFolderIDJson, _ := json.Marshal(*query.FolderId)
		params["new_folder_id"] = newFolderIDJson
	}
	err = r.esClient.UpdateByQuery(
		types.CredentialIndex,
		updateQuery,
		script,
		&params,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTag deletes a tag by ID from Elasticsearch
func (r *ESTagRepository) Delete(query DeleteTagQuery) error {
	if query.ID == "" {
		return fmt.Errorf("ID is required")
	}
	err := r.esClient.DeleteDocument(types.TagIndex, query.ID)
	if err != nil {
		return err
	}

	// Remove all the occurence of the tag inside the credential index

	script := fmt.Sprintf(`
		if (ctx._source.tags != null) {
			ctx._source.tags.removeIf(t -> t.id == '%s');
		}
	`, query.ID)
	updateQuery := esTypes.Query{
		Nested: &esTypes.NestedQuery{
			Path: "tags",
			Query: esTypes.Query{
				Term: map[string]esTypes.TermQuery{
					"tags.id": {Value: query.ID},
				},
			},
		},
	}

	err = r.esClient.UpdateByQuery(
		types.CredentialIndex,
		updateQuery,
		script,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchTags searches for tags in Elasticsearch based on the query
func (r *ESTagRepository) Search(query SearchTagQuery) (*SearchTagResult, error) {
	// Default limit and offset if not provided
	limit := 10
	if query.Limit != nil {
		limit = *query.Limit
	}
	offset := 0
	if query.Offset != nil {
		offset = *query.Offset
	}

	filters := []esTypes.Query{{
		Terms: &esTypes.TermsQuery{
			TermsQuery: map[string]esTypes.TermsQueryField{
				"folder_id": *query.FoldersScope,
			},
		},
	}}

	// Construct the search query
	res, total, err := r.esClient.Search(
		types.TagIndex,
		query.Name,
		[]string{"name^2"}, // Boost the name field
		filters,
	)
	if err != nil {
		return nil, err
	}

	// Parse the search results
	tags := make([]types.Tag, *total)
	for i, hit := range res {
		if err := json.Unmarshal(hit, &tags[i]); err != nil {
			return nil, fmt.Errorf("error unmarshalling hit source: %w", err)
		}
	}

	return &SearchTagResult{
		Tags:   tags,
		Total:  *total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
