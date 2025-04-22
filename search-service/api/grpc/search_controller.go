package grpc

import (
	"context"
	"errors"

	// Search service imports
	credentialService "github.com/DO-2K23-26/polypass-microservices/search-service/services/credential"
	folderService "github.com/DO-2K23-26/polypass-microservices/search-service/services/folder"
	tagService "github.com/DO-2K23-26/polypass-microservices/search-service/services/tags"

	// Following below is the imported generated gRPC code based on the protobuf definitions (.proto file)
	proto "github.com/DO-2K23-26/polypass-microservices/search-service/gen/search/api"
)

/* This is the gRPC API for the search service.
 * It defines an interface to search for folders, tags, and credentials.
 * It handles converting the incoming and outgoing gRPC responses.
 * Under the hood, it uses the respective services to perform the searches.
 * Search logic is implemented in the respective services.
 */

/*
 * Interfaces defined:
 * - ISearchGrpcApi
 *
 * Structs defined:
 * - SearchGrpcApi
 *
 * Functions implemented:
 * - NewSearchGrpcApi
 * - SearchFolders
 * - SearchTags
 * - SearchCredentials
*/

// SearchGrpcApi's interface
type ISearchGrpcApi interface {
	// API object constructor
	NewSearchGrpcApi(
		credentialService *credentialService.CredentialService,
		folderService *folderService.FolderService,
		tagService *tagService.TagService,
	) *SearchGrpcApi

	// Search endpoint handlers
	SearchFolders(ctx context.Context, req *proto.SearchFoldersRequest) (*proto.SearchFoldersResponse, error)
	SearchTags(ctx context.Context, req *proto.SearchTagsRequest) (*proto.SearchTagsResponse, error)
	SearchCredentials(ctx context.Context, req *proto.SearchCredentialsRequest) (*proto.SearchCredentialsResponse, error)
}

// SearchGrpcApi implements the ISearchGrpcApi interface
type SearchGrpcApi struct {
	//api.UnimplementedSearchServiceServer
	credentialService *credentialService.CredentialService
	folderService     *folderService.FolderService
	tagService        *tagService.TagService
}

// SearchGrpcApi constructor
func NewSearchGrpcApi(
	credentialService *credentialService.CredentialService,
	folderService *folderService.FolderService,
	tagService *tagService.TagService,
) *SearchGrpcApi {
	return &SearchGrpcApi{
		credentialService: credentialService,
		folderService:     folderService,
		tagService:        tagService,
	}
}

// SearchFolders searches for folders
func (api *SearchGrpcApi) SearchFolders(ctx context.Context, req *proto.SearchFoldersRequest) (*proto.SearchFoldersResponse, error) {
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
	result, err := api.folderService.SearchFolders(serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert service response to gRPC response
	response := &proto.SearchFoldersResponse{
		Total:  int32(result.Total),
		Limit:  int32(result.Limit),
		Offset: int32(result.Offset),
	}

	// Map folders
	response.Folders = make([]*proto.Folder, len(result.Folders))
	for i, folder := range result.Folders {
		response.Folders[i] = &proto.Folder{
			Id:   folder.ID,
			Name: folder.Name,
		}
	}

	return response, nil
}

// SearchTags searches for tags
func (api *SearchGrpcApi) SearchTags(ctx context.Context, req *proto.SearchTagsRequest) (*proto.SearchTagsResponse, error) {
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
	result, err := api.tagService.SearchTags(serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert service response to gRPC response
	response := &proto.SearchTagsResponse{
		Total:  int32(result.Total),
		Limit:  int32(result.Limit),
		Offset: int32(result.Offset),
	}

	// Map tags
	response.Tags = make([]*proto.Tag, len(result.Tags))
	for i, tag := range result.Tags {
		response.Tags[i] = &proto.Tag{
			Id:   tag.ID,
			Name: tag.Name,
		}
	}

	return response, nil
}

// SearchCredentials searches for credentials
func (api *SearchGrpcApi) SearchCredentials(ctx context.Context, req *proto.SearchCredentialsRequest) (*proto.SearchCredentialsResponse, error) {
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

	serviceReq := credentialService.SearchCredentialsRequest{
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
	result, err := api.credentialService.SearchCredentials(serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert service response to gRPC response
	response := &proto.SearchCredentialsResponse{
		Total:  int32(result.Total),
		Limit:  int32(result.Limit),
		Offset: int32(result.Offset),
	}

	// Map credentials
	response.Credentials = make([]*proto.Credential, len(result.Credentials))
	for i, cred := range result.Credentials {
		// Create credential
		credential := &proto.Credential{
			Id:       cred.ID,
			Title:    cred.Title,
			FolderId: cred.FolderID,
		}

		// Add folder if present
		if cred.Folder != nil {
			credential.Folder = &proto.Folder{
				Id:   cred.Folder.ID,
				Name: cred.Folder.Name,
			}
		}

		// Add tags if present
		if len(cred.Tags) > 0 {
			credential.Tags = make([]*proto.Tag, len(cred.Tags))
			for j, tag := range cred.Tags {
				credential.Tags[j] = &proto.Tag{
					Id:   tag.ID,
					Name: tag.Name,
				}
			}
		}

		response.Credentials[i] = credential
	}

	return response, nil
}
