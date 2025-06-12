// filepath: gateway/core/search.go
package core

import (
	"context"
	"github.com/DO-2K23-26/polypass-microservices/gateway/proto/search"
)

// SearchService defines methods to search folders, tags, and credentials via gRPC
type SearchService interface {
	SearchFolders(ctx context.Context, req *search.SearchFoldersRequest) ([]*search.Folder, int32, error)
	SearchTags(ctx context.Context, req *search.SearchTagsRequest) ([]*search.Tag, int32, error)
	SearchCredentials(ctx context.Context, req *search.SearchCredentialsRequest) ([]*search.Credential, int32, error)
}

// searchService implements SearchService
type searchService struct {
	client search.SearchServiceClient
}

// NewSearchService creates a new SearchService
func NewSearchService(client search.SearchServiceClient) SearchService {
	return &searchService{client: client}
}

func (s *searchService) SearchFolders(ctx context.Context, req *search.SearchFoldersRequest) ([]*search.Folder, int32, error) {
	resp, err := s.client.SearchFolders(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return resp.Folders, resp.Total, nil
}

func (s *searchService) SearchTags(ctx context.Context, req *search.SearchTagsRequest) ([]*search.Tag, int32, error) {
	resp, err := s.client.SearchTags(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return resp.Tags, resp.Total, nil
}

func (s *searchService) SearchCredentials(ctx context.Context, req *search.SearchCredentialsRequest) ([]*search.Credential, int32, error) {
	resp, err := s.client.SearchCredentials(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return resp.Credentials, resp.Total, nil
}
