package tag

import (
	"errors"

	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/tags"
	"github.com/DO-2K23-26/polypass-microservices/search-service/services/folder"
)

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrTagNotFound       = errors.New("tag not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotAuthorized = errors.New("user not authorized to access this folder")
)

type TagService struct {
	tagRepository tags.ITagRepository
	folderService folder.FolderService
}

func NewTagService(tagRepo tags.ITagRepository) *TagService {
	return &TagService{
		tagRepository: tagRepo,
	}
}

// CreateTag creates a new tag
func (s *TagService) Create(req CreateTagRequest) (*CreateTagResponse, error) {
	if req.Name == "" {
		return nil, ErrInvalidRequest
	}

	result, err := s.tagRepository.Create(tags.CreateTagQuery{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &CreateTagResponse{
		Tag: *result.Tag,
	}, nil
}

// GetTag retrieves a tag by ID
func (s *TagService) Get(req GetTagRequest) (*CreateTagResponse, error) {
	if req.ID == "" {
		return nil, ErrInvalidRequest
	}

	result, err := s.tagRepository.Get(tags.GetTagQuery{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	return &CreateTagResponse{
		Tag: *result.Tag,
	}, nil
}

func (s *TagService) MGet(req MGetTagRequest) (*MGetTagResponse, error) {
	if len(req.IDs) == 0 {
		return nil, ErrInvalidRequest
	}

	result, err := s.tagRepository.MGet(tags.MGetTagQuery{
		IDs: req.IDs,
	})
	if err != nil {
		return nil, err
	}

	return &MGetTagResponse{
		Tags: result.Tags,
	}, nil

}

// UpdateTag updates an existing tag
func (s *TagService) Update(req UpdateTagRequest) error {
	if req.ID == "" {
		return ErrInvalidRequest
	}

	err := s.tagRepository.Update(tags.UpdateTagQuery{})
	if err != nil {
		return err
	}

	return nil
}

// DeleteTag deletes a tag by ID
func (s *TagService) Delete(req DeleteTagRequest) error {
	if req.ID == "" {
		return ErrInvalidRequest
	}

	return s.tagRepository.Delete(tags.DeleteTagQuery{
		ID: req.ID,
	})
}

// SearchTags searches for tags based on criteria
func (s *TagService) Search(req SearchTagsRequest) (*SearchTagsResponse, error) {
	// Get user to determine folder access scope

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

	//Extract folder IDs from the response

	folderIDs := make([]string, len(res.Folders))
	for i, folder := range res.Folders {
		folderIDs[i] = folder.ID
	}
	// Perform the search
	searchResult, err := s.tagRepository.Search(tags.SearchTagQuery{
		Name:         req.SearchQuery,
		Limit:        &limit,
		Offset:       &offset,
		FoldersScope: &folderIDs,
	})
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	response := &SearchTagsResponse{
		Tags:   searchResult.Tags,
		Limit:  searchResult.Limit,
		Offset: searchResult.Offset,
		Total:  searchResult.Total,
	}

	return response, nil
}
