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
	credentialRepository credential.ICredentialRepository
	folderService        folder.FolderService
	userService          user.IUserRepository
	tagService           tag.TagService
}

func NewCredentialService(
	credentialRepository credential.ICredentialRepository,
) *CredentialService {
	return &CredentialService{
		credentialRepository: credentialRepository,
	}
}

// CreateCredential creates a new credential
func (s *CredentialService) Create(req CreateCredentialRequest) (*CreateCredentialResponse, error) {
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
	result, err := s.credentialRepository.Create(credential.CreateCredentialQuery{
		ID:     req.ID,
		Title:  req.Title,
		Folder: res.Folder,
	})
	if err != nil {
		return nil, err
	}

	return &CreateCredentialResponse{
		Credential: result.Credential,
	}, nil
}

// GetCredential retrieves a credential by ID
func (s *CredentialService) Get(req GetCredentialRequest) (*GetCredentialResponse, error) {
	if req.ID == "" {
		return nil, ErrInvalidRequest
	}

	// Get the credential
	result, err := s.credentialRepository.Get(credential.GetCredentialQuery{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrCredentialNotFound
	}

	response := result.Credential
	return &GetCredentialResponse{
		Credential: response,
	}, nil
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
	return s.credentialRepository.Update(credential.UpdateCredentialQuery{
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
	return s.credentialRepository.Delete(credential.DeleteCredentialQuery{
		ID: req.ID,
	})
}

// SearchCredentials searches for credentials based on criteria
func (s *CredentialService) Search(req SearchCredentialsRequest) (*SearchCredentialsResponse, error) {
	if req.UserID == "" {
		return nil, ErrInvalidRequest
	}

	// Set default limit and offset if not provided
	limit := 10
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	offset := 0
	if req.Page != nil && *req.Page >= 0 {
		offset = *req.Page * limit
	}

	res, err := s.folderService.GetFromUser(folder.GetUserFoldersRequest{UserID: req.UserID})
	if err != nil {
		return nil, err
	}

	var foldersScope []string
	for _, folder := range res.Folders {
		foldersScope = append(foldersScope, folder.ID)
	}

	// Perform the search
	searchResult, err := s.credentialRepository.Search(credential.SearchCredentialQuery{
		TagIds:       req.TagIDs,
		Limit:        &limit,
		Offset:       &offset,
		FoldersScope: foldersScope,
	})
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	response := &SearchCredentialsResponse{
		Credentials: searchResult.Credentials,
		Total:       searchResult.Total,
	}

	return response, nil
}

// AddTagsToCredential adds tags to a credential (this would need a new repository method)
func (s *CredentialService) AddTags(credentialID string, tagIDs []string) error {
	// Validate input
	if credentialID == "" {
		return ErrInvalidRequest
	}

	res, err := s.tagService.MGet(tag.MGetTagRequest{IDs: tagIDs})
	if err != nil {
		return err
	}
	// Call the repository method to add tags
	return s.credentialRepository.AddTags(credential.AddTagsToCredentialQuery{
		ID:   credentialID,
		Tags: res.Tags,
	})

}

// RemoveTagsFromCredential removes tags from a credential (this would need a new repository method)
func (s *CredentialService) RemoveTags(req RemoveTagsFromCredentialRequest) error {
	if req.ID == "" {
		return ErrInvalidRequest
	}
	return s.credentialRepository.RemoveTags(credential.RemoveTagsFromCredentialQuery{
		ID:     req.ID,
		TagIds: req.TagIds,
	})
}
