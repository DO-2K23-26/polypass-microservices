package grpc

import (
	"context"
	"errors"

	"github.com/DO-2K23-26/polypass-microservices/search-service/gen/search/api"
	"github.com/DO-2K23-26/polypass-microservices/search-service/services/credential"
	folderService "github.com/DO-2K23-26/polypass-microservices/search-service/services/folder"
	tagService "github.com/DO-2K23-26/polypass-microservices/search-service/services/tags"
)

// SearchServiceServer implements the gRPC search service
type SearchServiceServer struct {
	api.UnimplementedSearchServiceServer
	credentialService *credential.CredentialService
	folderService     *folderService.FolderService
	tagService        *tagService.TagService
}

// NewSearchServiceServer creates a new search service server
func NewSearchServiceServer(
	credentialService *credential.CredentialService,
	folderService *folderService.FolderService,
	tagService *tagService.TagService,
) *SearchServiceServer {
	return &SearchServiceServer{
		credentialService: credentialService,
		folderService:     folderService,
		tagService:        tagService,
	}
}

// SearchFolders searches for folders
func (s *SearchServiceServer) SearchFolders(ctx context.Context, req *api.SearchFoldersRequest) (*api.SearchFoldersResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, errors.New("user ID is required")
	}

	// Convert gRPC request to service request
	limit := int(req.Limit)
	offset := int(req.Offset)
	serviceReq := folderService.SearchFoldersRequest{
		Name:   req.Name,
		UserID: req.UserId,
		Limit:  &limit,
		Offset: &offset,
	}

	// Call service
	result, err := s.folderService.SearchFolders(serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert service response to gRPC response
	response := &api.SearchFoldersResponse{
		Total:  int32(result.Total),
		Limit:  int32(result.Limit),
		Offset: int32(result.Offset),
	}

	// Map folders
	response.Folders = make([]*api.Folder, len(result.Folders))
	for i, folder := range result.Folders {
		response.Folders[i] = &api.Folder{
			Id:   folder.ID,
			Name: folder.Name,
		}
	}

	return response, nil
}

// SearchTags searches for tags
func (s *SearchServiceServer) SearchTags(ctx context.Context, req *api.SearchTagsRequest) (*api.SearchTagsResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, errors.New("user ID is required")
	}

	// Convert gRPC request to service request
	limit := int(req.Limit)
	offset := int(req.Offset)
	var folderID *string
	if req.FolderId != "" {
		folderID = &req.FolderId
	}

	serviceReq := tagService.SearchTagsRequest{
		Name:     req.Name,
		FolderID: folderID,
		Limit:    &limit,
		Offset:   &offset,
		UserID:   req.UserId,
	}

	// Call service
	result, err := s.tagService.SearchTags(serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert service response to gRPC response
	response := &api.SearchTagsResponse{
		Total:  int32(result.Total),
		Limit:  int32(result.Limit),
		Offset: int32(result.Offset),
	}

	// Map tags
	response.Tags = make([]*api.Tag, len(result.Tags))
	for i, tag := range result.Tags {
		response.Tags[i] = &api.Tag{
			Id:   tag.ID,
			Name: tag.Name,
		}
	}

	return response, nil
}

// SearchCredentials searches for credentials
func (s *SearchServiceServer) SearchCredentials(ctx context.Context, req *api.SearchCredentialsRequest) (*api.SearchCredentialsResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, errors.New("user ID is required")
	}

	// Convert gRPC request to service request
	limit := int(req.Limit)
	offset := int(req.Offset)

	var folderID *string
	if req.FolderId != "" {
		folderID = &req.FolderId
	}

	var folderName *string
	if req.FolderName != "" {
		folderName = &req.FolderName
	}

	var tagName *string
	if req.TagName != "" {
		tagName = &req.TagName
	}

	var tagIDs *[]string
	if len(req.TagIds) > 0 {
		tags := req.TagIds
		tagIDs = &tags
	}

	serviceReq := credential.SearchCredentialsRequest{
		Title:      req.Title,
		FolderID:   folderID,
		FolderName: folderName,
		TagIDs:     tagIDs,
		TagName:    tagName,
		Limit:      &limit,
		Offset:     &offset,
		UserID:     req.UserId,
	}

	// Call service
	result, err := s.credentialService.SearchCredentials(serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert service response to gRPC response
	response := &api.SearchCredentialsResponse{
		Total:  int32(result.Total),
		Limit:  int32(result.Limit),
		Offset: int32(result.Offset),
	}

	// Map credentials
	response.Credentials = make([]*api.Credential, len(result.Credentials))
	for i, cred := range result.Credentials {
		// Create credential
		credential := &api.Credential{
			Id:       cred.ID,
			Title:    cred.Title,
			FolderId: cred.FolderID,
		}

		// Add folder if present
		if cred.Folder != nil {
			credential.Folder = &api.Folder{
				Id:   cred.Folder.ID,
				Name: cred.Folder.Name,
			}
		}

		// Add tags if present
		if len(cred.Tags) > 0 {
			credential.Tags = make([]*api.Tag, len(cred.Tags))
			for j, tag := range cred.Tags {
				credential.Tags[j] = &api.Tag{
					Id:   tag.ID,
					Name: tag.Name,
				}
			}
		}

		response.Credentials[i] = credential
	}

	return response, nil
}
