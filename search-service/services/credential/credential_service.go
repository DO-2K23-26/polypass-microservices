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
	if req.Name == nil || *req.Name == "" {
		return nil, ErrInvalidRequest
	}

	if req.FolderID == nil || *req.FolderID == "" {
		return nil, ErrInvalidRequest
	}

	var tagIDs []string
	if req.TagIDs != nil {
		tagIDs = *req.TagIDs
	}

	result, err := s.credentialRepository.Create(credential.CreateCredentialQuery{
		ID:        req.CredentialID,
		Name:      *req.Name,
		FolderID:  *req.FolderID,
		TagIDs:    tagIDs,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	return &CreateCredentialResponse{
		Credential: types.Credential{
			ID:        result.Credential.ID,
			Name:      result.Credential.Name,
			FolderID:  result.Credential.FolderID,
			TagIDs:    result.Credential.TagIDs,
			CreatedAt: result.Credential.CreatedAt,
			UpdatedAt: result.Credential.UpdatedAt,
		},
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
	if req.CredentialID == "" {
		return ErrInvalidRequest
	}
	query := credential.UpdateCredentialQuery{
		ID:       req.CredentialID,
		Name:    req.Name,
		FolderID: req.FolderID,
		TagIDs:   req.TagIDs,
	}
	return s.credentialRepository.Update(query)
}

// DeleteCredential deletes a credential by ID
func (s *CredentialService) Delete(req DeleteCredentialRequest) error {
	if req.CredentialID == "" {
		return ErrInvalidRequest
	}
	query := credential.DeleteCredentialQuery{
		ID: req.CredentialID,
	}
	return s.credentialRepository.Delete(query)
}

// SearchCredentials searches for credentials based on criteria
func (s *CredentialService) Search(req SearchCredentialsRequest) (*SearchCredentialsResponse, error) {
	if req.UserID == "" {
		return nil, ErrInvalidRequest
	}
	// Ici, il faudrait idéalement remplir FoldersScope avec les dossiers accessibles à l'utilisateur
	query := credential.SearchCredentialQuery{
		SearchQuery:  req.SearchQuery,
		TagIds:       req.TagIDs,
		Limit:        req.Limit,
		Offset:       req.Page,
		FoldersScope: nil, // à adapter si besoin
	}
	result, err := s.credentialRepository.Search(query)
	if err != nil {
		return nil, err
	}
	return &SearchCredentialsResponse{
		Credentials: result.Credentials,
		Total:       result.Total,
	}, nil
}

// AddTagsToCredential adds tags to a credential
func (s *CredentialService) AddTags(req AddTagsToCredentialRequest) error {
	if req.ID == "" {
		return ErrInvalidRequest
	}
	if len(req.TagIds) == 0 {
		return ErrInvalidRequest
	}
	query := credential.AddTagsToCredentialQuery{
		ID:   req.ID,
		Tags: nil, // à adapter selon la logique de ton repo
	}
	return s.credentialRepository.AddTags(query)
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
