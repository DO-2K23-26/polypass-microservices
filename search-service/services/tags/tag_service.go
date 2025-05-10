package tag

import (
	"errors"
	"strings"

	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/tags"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"
)

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrTagNotFound       = errors.New("tag not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotAuthorized = errors.New("user not authorized to access this folder")
)

type TagService struct {
	tagRepo tags.ITagRepository
}

func NewTagService(tagRepo tags.ITagRepository, userRepo user.IUserRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// CreateTag creates a new tag
func (s *TagService) Create(req CreateTagRequest) (*TagResponse, error) {
	if req.Name == "" {
		return nil, ErrInvalidRequest
	}

	result, err := s.tagRepo.Create(tags.CreateTagQuery{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &TagResponse{
		ID:   result.Tag.ID,
		Name: result.Tag.Name,
	}, nil
}

// GetTag retrieves a tag by ID
func (s *TagService) Get(req GetTagRequest) (*TagResponse, error) {
	if req.ID == "" {
		return nil, ErrInvalidRequest
	}

	result, err := s.tagRepo.Get(tags.GetTagQuery{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	return &TagResponse{
		ID:   result.Tag.ID,
		Name: result.Tag.Name,
	}, nil
}

func (s *TagService) MGet(req MGetTagRequest) (*MGetTagResponse, error) {
	if len(req.IDs) == 0 {
		return nil, ErrInvalidRequest
	}

	result, err := s.tagRepo.MGet(tags.MGetTagQuery{
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
func (s *TagService) Update(req UpdateTagRequest) (*TagResponse, error) {
	if req.ID == "" || req.Name == "" {
		return nil, ErrInvalidRequest
	}

	// Normalize tag name (e.g., lowercase, trim spaces)
	normalizedName := strings.TrimSpace(strings.ToLower(req.Name))
	if normalizedName == "" {
		return nil, ErrInvalidRequest
	}

	// result, err := s.tagRepo.Update(tags.UpdateTagQuery{
	// 	ID:   req.ID,
	// 	Name: normalizedName,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

// DeleteTag deletes a tag by ID
func (s *TagService) Delete(req DeleteTagRequest) error {
	if req.ID == "" {
		return ErrInvalidRequest
	}

	return s.tagRepo.Delete(tags.DeleteTagQuery{
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

	// Perform the search
	searchResult, err := s.tagRepo.Search(tags.SearchTagQuery{
		Name:         req.SearchQuery,
		Limit:        &limit,
		Offset:       &offset,
		FoldersScope: &req.FolderIDs,
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
