package grpc

import (
	"context"
	"errors"

	types "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	api "github.com/DO-2K23-26/polypass-microservices/search-service/gen/search/api"
	credentialService "github.com/DO-2K23-26/polypass-microservices/search-service/services/credential"
	folderService "github.com/DO-2K23-26/polypass-microservices/search-service/services/folder"
	tagService "github.com/DO-2K23-26/polypass-microservices/search-service/services/tags"
)

// SearchServiceServer implements the gRPC search service
type SearchServiceServer struct {
	api.UnimplementedSearchServiceServer
	credentialService *credentialService.CredentialService
	folderService     *folderService.FolderService
	tagService        *tagService.TagService
}

// NewSearchServiceServer creates a new search service server
func NewSearchServiceServer(
	credentialService *credentialService.CredentialService,
	folderService *folderService.FolderService,
	tagService *tagService.TagService,
) *SearchServiceServer {
	if credentialService == nil || folderService == nil || tagService == nil {
		panic("Services must not be nil")
	}
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
	page := int(req.Page)
	serviceReq := folderService.SearchFoldersRequest{
		SearchQuery: req.SearchQuery,
		UserID:      req.UserId,
		Limit:       &limit,
		Page:        &page,
	}

	// Call service
	result, err := s.folderService.Search(serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert service response to gRPC response
	response := &api.SearchFoldersResponse{
		Total: int32(result.Total),
	}

	// Map folders
	response.Folders = make([]*api.Folder, len(result.Folders))
	for i, folder := range result.Folders {
		response.Folders[i] = &api.Folder{
			Id:       folder.ID,
			Name:     folder.Name,
			ParentId: folder.ParentID.String,
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
	page := int(req.Page)

	serviceReq := tagService.SearchTagsRequest{
		SearchQuery: req.SearchQuery,
		Limit:       &limit,
		Page:        &page,
		UserID:      req.UserId,
	}

	// Call service
	result, err := s.tagService.Search(serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert service response to gRPC response
	response := &api.SearchTagsResponse{
		Total: int32(result.Total),
	}

	// Map tags
	response.Tags = make([]*api.Tag, len(result.Tags))
	for i, tag := range result.Tags {
		response.Tags[i] = &api.Tag{
			Id:       tag.ID,
			Name:     tag.Name,
			FolderId: tag.FolderId,
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
	page := int(req.Page)

	var folderID *string
	if req.FolderId != "" {
		folderID = &req.FolderId
	}

	var tagIDs *[]string
	if len(req.TagIds) > 0 {
		tagIDs = &req.TagIds
	}

	serviceReq := credentialService.SearchCredentialsRequest{
		SearchQuery: req.SearchQuery,
		FolderID:    folderID,
		TagIDs:      tagIDs,
		Page:        &page,
		Limit:       &limit,
		UserID:      req.UserId,
	}

	// Call service
	response, err := s.credentialService.Search(serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert response to proto
	protoResponse := &api.SearchCredentialsResponse{
		Credentials: make([]*api.Credential, len(response.Credentials)),
		Total:       int32(response.Total),
	}

	for i, cred := range response.Credentials {
		protoResponse.Credentials[i] = convertCredentialToProto(cred)
	}

	return protoResponse, nil
}

// Convert credential to proto
func convertCredentialToProto(cred types.Credential) *api.Credential {
	return &api.Credential{
		Id:       cred.ID,
		Title:    cred.Name,
		FolderId: cred.FolderID,
	}
}
