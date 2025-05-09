package credential

import (
	"errors"

	"github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/credential"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"
	"github.com/DO-2K23-26/polypass-microservices/search-service/services/folder"
	tag "github.com/DO-2K23-26/polypass-microservices/search-service/services/tags"
)

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrCredentialNotFound = errors.New("credential not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotAuthorized  = errors.New("user not authorized to access this credential")
	ErrFolderNotFound     = errors.New("folder not found")
	ErrTagNotFound        = errors.New("one or more tags not found")
)

type CredentialService struct {
	credentialRepo credential.ICredentialRepository
	folderService  folder.FolderService
	userService    user.IUserRepository
	tagService     tag.TagService
}

func NewCredentialService(
	credentialRepo credential.ICredentialRepository,
	userRepo user.IUserRepository,
) *CredentialService {
	return &CredentialService{
		credentialRepo: credentialRepo,
	}
}

// CreateCredential creates a new credential
func (s *CredentialService) Create(req CreateCredentialRequest) (*CredentialResponse, error) {
	// Validate required fields
	if req.Title == "" || req.FolderID == "" || req.ID == "" {
		return nil, ErrInvalidRequest
	}

	res, err := s.folderService.Get(folder.GetFolderRequest{
		ID: req.FolderID,
	})
	if err != nil {
		return nil, err
	}
	// Create the credential
	result, err := s.credentialRepo.Create(credential.CreateCredentialQuery{
		ID:     req.ID,
		Title:  req.Title,
		Folder: res.Folder,
	})
	if err != nil {
		return nil, err
	}

	return &CredentialResponse{
		ID:     result.Credential.ID,
		Title:  result.Credential.Title,
		Tags:   ConvertToTagResponses(result.Credential.Tags),
		Folder: ConvertToFolderResponse(result.Credential.Folder),
	}, nil
}

// GetCredential retrieves a credential by ID
func (s *CredentialService) Get(req GetCredentialRequest) (*CredentialResponse, error) {
	if req.ID == "" {
		return nil, ErrInvalidRequest
	}

	// Get the credential
	result, err := s.credentialRepo.Get(credential.GetCredentialQuery{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrCredentialNotFound
	}

	response := ConvertToCredentialResponse(result.Credential)
	return &response, nil
}

// UpdateCredential updates an existing credential
func (s *CredentialService) Update(req UpdateCredentialRequest) error {
	if req.ID == "" {
		return ErrInvalidRequest
	}
	var newFolder *types.Folder
	if req.FolderID != "" {
		res, err := s.folderService.Get(folder.GetFolderRequest{ID: req.ID})
		if err != nil {
			return err
		}
		newFolder = &res.Folder
	}

	// Update the credential
	return s.credentialRepo.Update(credential.UpdateCredentialQuery{
		ID:     req.ID,
		Title:  &req.Title,
		Folder: newFolder,
	})
}

// DeleteCredential deletes a credential by ID
func (s *CredentialService) Delete(req DeleteCredentialRequest) error {
	if req.ID == "" {
		return ErrInvalidRequest
	}

	// Delete the credential
	return s.credentialRepo.Delete(credential.DeleteCredentialQuery{
		ID: req.ID,
	})
}

// SearchCredentials searches for credentials based on criteria
func (s *CredentialService) Search(req SearchCredentialsRequest) (*SearchCredentialsResponse, error) {
	if req.UserID == "" {
		return nil, ErrInvalidRequest
	}

	// Set default limit and offset if not provided
	// limit := 10
	// if req.Limit != nil && *req.Limit > 0 {
	// 	limit = *req.Limit
	// }

	// offset := 0
	// if req.Offset != nil && *req.Offset >= 0 {
	// 	offset = *req.Offset
	// }

	// Perform the search
	// searchResult, err := s.credentialRepo.SearchCredentials(credential.SearchCredentialQuery{
	// 	Title:        req.Title,
	// 	FolderId:     req.FolderID,
	// 	FolderName:   req.FolderName,
	// 	TagIds:       req.TagIDs,
	// 	TagName:      req.TagName,
	// 	Limit:        &limit,
	// 	Offset:       &offset,
	// 	FoldersScope: &folderIds,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// Convert to response DTO
	response := &SearchCredentialsResponse{
		// Credentials: ConvertToCredentialResponses(searchResult.Credentials),
		// Total:       searchResult.Total,
		// Limit:       searchResult.Limit,
		// Offset:      searchResult.Offset,
	}

	return response, nil
}

// AddTagsToCredential adds tags to a credential (this would need a new repository method)
func (s *CredentialService) AddTags(credentialID string, tagIDs []string) error {
	// Validate input
	if credentialID == "" {
		return ErrInvalidRequest
	}

	res, err := s.tagService.MGet(tag.MGetTagRequest{})
	if err != nil {
		return err
	}
	// Call the repository method to add tags
	return s.credentialRepo.AddTags(credential.AddTagsToCredentialQuery{
		ID:   credentialID,
		Tags: res.Tags,
	})

}

// RemoveTagsFromCredential removes tags from a credential (this would need a new repository method)
func (s *CredentialService) RemoveTags(req RemoveTagsFromCredentialRequest) error {
	if req.ID == "" {
		return ErrInvalidRequest
	}
	return s.credentialRepo.RemoveTags(credential.RemoveTagsFromCredentialQuery{
		ID:     req.ID,
		TagIds: req.TagIds,
	})
}
