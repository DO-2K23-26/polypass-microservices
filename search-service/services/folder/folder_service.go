package folder

import (
	"errors"

	"slices"

	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/folder"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"
)

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrFolderNotFound    = errors.New("folder not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotAuthorized = errors.New("user not authorized to access this folder")
)

type FolderService struct {
	folderRepo folder.IFolderRepository
	userRepo   user.IUserRepository
}

func NewFolderService(folderRepo folder.IFolderRepository, userRepo user.IUserRepository) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
		userRepo:   userRepo,
	}
}

// CreateFolder creates a new folder
func (s *FolderService) CreateFolder(req CreateFolderRequest) (*FolderResponse, error) {
	if req.Name == "" {
		return nil, ErrInvalidRequest
	}

	result, err := s.folderRepo.CreateFolder(folder.CreateFolderQuery{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &FolderResponse{
		ID:   result.Folder.ID,
		Name: result.Folder.Name,
	}, nil
}

// GetFolder retrieves a folder by ID, checking user permissions
func (s *FolderService) GetFolder(req GetFolderRequest) (*FolderResponse, error) {
	if req.ID == "" || req.UserID == "" {
		return nil, ErrInvalidRequest
	}

	// Check if user has access to the folder
	userResult, err := s.userRepo.Get(user.GetUserQuery{
		ID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	if userResult == nil || userResult.User.ID == "" {
		return nil, ErrUserNotFound
	}

	folderIds := make([]string, len(userResult.User.Folders))
	
	for _, folder := range userResult.User.Folders {
		folderIds = append(folderIds, folder.ID)
	}
	// Check if user has access to the folder
	hasAccess := slices.Contains(folderIds, req.ID)
	if !hasAccess {
		return nil, ErrUserNotAuthorized
	}

	// Get the folder
	result, err := s.folderRepo.GetFolder(folder.GetFolderQuery{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrFolderNotFound
	}

	return &FolderResponse{
		ID:   result.Folder.ID,
		Name: result.Folder.Name,
	}, nil
}

// UpdateFolder updates an existing folder
func (s *FolderService) UpdateFolder(req UpdateFolderRequest) (*FolderResponse, error) {
	if req.ID == "" || req.Name == "" {
		return nil, ErrInvalidRequest
	}

	result, err := s.folderRepo.UpdateFolder(folder.UpdateFolderQuery{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &FolderResponse{
		ID:   result.Folder.ID,
		Name: result.Folder.Name,
	}, nil
}

// DeleteFolder deletes a folder by ID, checking user permissions
func (s *FolderService) DeleteFolder(req DeleteFolderRequest) error {
	if req.ID == "" || req.UserID == "" {
		return ErrInvalidRequest
	}

	// Check if user has access to the folder
	userResult, err := s.userRepo.Get(user.GetUserQuery{
		ID: req.UserID,
	})
	if err != nil {
		return err
	}
	if userResult == nil || userResult.User.ID == "" {
		return ErrUserNotFound
	}
	
	folderIds := make([]string, len(userResult.User.Folders))
	
	for _, folder := range userResult.User.Folders {
		folderIds = append(folderIds, folder.ID)
	}
	// Check if user has access to the folder
	hasAccess := slices.Contains(folderIds, req.ID)
	if !hasAccess {
		return ErrUserNotAuthorized
	}

	return s.folderRepo.DeleteFolder(folder.DeleteFolderQuery{
		ID: req.ID,
	})
}

// SearchFolders searches for folders based on criteria
func (s *FolderService) SearchFolders(req SearchFoldersRequest) (*SearchFoldersResponse, error) {
	if req.UserID == "" {
		return nil, ErrInvalidRequest
	}

	// Get user to determine folder access scope
	userResult, err := s.userRepo.Get(user.GetUserQuery{
		ID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	if userResult == nil || userResult.User.ID == "" {
		return nil, ErrUserNotFound
	}

	// Set default limit and offset if not provided
	limit := 10
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	offset := 0
	if req.Offset != nil && *req.Offset >= 0 {
		offset = *req.Offset
	}

	folderIds := make([]string, len(userResult.User.Folders))
	
	for _, folder := range userResult.User.Folders {
		folderIds = append(folderIds, folder.ID)
	}
	// Perform the search with user's folder access scope
	searchResult, err := s.folderRepo.SearchFolder(folder.SearchFolderQuery{
		ID:           req.ID,
		Name:         req.Name,
		Limit:        &limit,
		Offset:       &offset,
		FoldersScope: &folderIds,
	})
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	response := &SearchFoldersResponse{
		Folders: ConvertToFoldersResponse(searchResult.Folders),
		Total:   searchResult.Total,
		Limit:   searchResult.Limit,
		Offset:  searchResult.Offset,
	}

	return response, nil
}
