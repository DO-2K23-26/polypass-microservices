package folder

import (
	"errors"

	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/folder"
	"github.com/DO-2K23-26/polypass-microservices/search-service/services/user"
)

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrFolderNotFound    = errors.New("folder not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotAuthorized = errors.New("user not authorized to access this folder")
)

type FolderService struct {
	folderRepo  folder.IFolderRepository
	userService user.UserService
}

func NewFolderService(folderRepo folder.IFolderRepository, userService user.UserService) *FolderService {
	return &FolderService{
		folderRepo:  folderRepo,
		userService: userService,
	}
}

// CreateFolder creates a new folder
func (s *FolderService) CreateFolder(req CreateFolderRequest) (*FolderResponse, error) {
	if req.Name == "" {
		return nil, ErrInvalidRequest
	}

	query := folder.CreateFolderQuery{
		ID:   req.ID,
		Name: req.Name,
	}

	if req.ParentID != nil {
		//TO DO: Check the right of the user to create a folder inside another folder
		query.ParentID = req.ParentID
	}

	result, err := s.folderRepo.Create(query)
	if err != nil {
		return nil, err
	}

	return &FolderResponse{
		Folder: result.Folder,
	}, nil
}

// GetFolder retrieves a folder by ID
func (s *FolderService) Get(req GetFolderRequest) (*FolderResponse, error) {
	if req.ID == "" {
		return nil, ErrInvalidRequest
	}

	// Get the folder
	result, err := s.folderRepo.Get(folder.GetFolderQuery{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	return &FolderResponse{
		Folder: result.Folder,
	}, nil
}

// UpdateFolder updates an existing folder
func (s *FolderService) Update(req UpdateFolderRequest) (*FolderResponse, error) {
	if req.ID == "" || req.Name == "" {
		return nil, ErrInvalidRequest
	}

	result, err := s.folderRepo.Update(folder.UpdateFolderQuery{
		ID:   req.ID,
		Name: &req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &FolderResponse{
		Folder: result.Folder,
	}, nil
}

// DeleteFolder deletes a folder by ID, checking user permissions
func (s *FolderService) Delete(req DeleteFolderRequest) error {
	if req.ID == "" || req.UserID == "" {
		return ErrInvalidRequest
	}

	return s.folderRepo.Delete(folder.DeleteFolderQuery{
		ID: req.ID,
	})
}

// SearchFolders searches for folders based on criteria
func (s *FolderService) Search(req SearchFoldersRequest) (*SearchFoldersResponse, error) {
	// Set default limit and offset if not provided
	limit := 10
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	offset := 0
	if req.Page != nil && *req.Page >= 0 {
		offset = *req.Page * limit
	}

	res, err := s.userService.GetFolders(user.GetFoldersRequest{UserID: req.UserID})
	if err != nil {
		return nil, err
	}

	folderIds := make([]string, len(res.Folders))
	for i, folder := range res.Folders {
		folderIds[i] = folder.ID
	}

	// Perform the search with user's folder access scope
	searchResult, err := s.folderRepo.Search(folder.SearchFolderQuery{
		Name:         req.SearchQuery,
		Limit:        &limit,
		Offset:       &offset,
		FoldersScope: &folderIds,
	})
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	response := &SearchFoldersResponse{
		Folders: searchResult.Folders,
		Total:   searchResult.Total,
	}

	return response, nil
}
