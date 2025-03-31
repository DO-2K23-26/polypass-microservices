package tag

import (
	"errors"
	"strings"

	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/tags"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"
	"slices"
)

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrTagNotFound       = errors.New("tag not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotAuthorized = errors.New("user not authorized to access this folder")
)

type TagService struct {
	tagRepo  tags.TagRepository
	userRepo user.UserRepository
}

func NewTagService(tagRepo tags.TagRepository, userRepo user.UserRepository) *TagService {
	return &TagService{
		tagRepo:  tagRepo,
		userRepo: userRepo,
	}
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(req CreateTagRequest) (*TagResponse, error) {
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
func (s *TagService) GetTag(req GetTagRequest) (*TagResponse, error) {
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

// UpdateTag updates an existing tag
func (s *TagService) UpdateTag(req UpdateTagRequest) (*TagResponse, error) {
	if req.ID == "" || req.Name == "" {
		return nil, ErrInvalidRequest
	}

	// Normalize tag name (e.g., lowercase, trim spaces)
	normalizedName := strings.TrimSpace(strings.ToLower(req.Name))
	if normalizedName == "" {
		return nil, ErrInvalidRequest
	}

	result, err := s.tagRepo.Update(tags.UpdateTagQuery{
		ID:   req.ID,
		Name: normalizedName,
	})
	if err != nil {
		return nil, err
	}

	return &TagResponse{
		ID:   result.Tag.ID,
		Name: result.Tag.Name,
	}, nil
}

// DeleteTag deletes a tag by ID
func (s *TagService) DeleteTag(req DeleteTagRequest) error {
	if req.ID == "" {
		return ErrInvalidRequest
	}

	return s.tagRepo.Delete(tags.DeleteTagQuery{
		ID: req.ID,
	})
}

// SearchTags searches for tags based on criteria
func (s *TagService) SearchTags(req SearchTagsRequest) (*SearchTagsResponse, error) {
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

	// If a specific folder ID is requested, verify the user has access to it
	if req.FolderID != nil && *req.FolderID != "" {
		hasAccess := slices.Contains(userResult.User.FolderIds, *req.FolderID)
		if !hasAccess {
			return nil, ErrUserNotAuthorized
		}
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

	// Perform the search
	searchResult, err := s.tagRepo.SearchTags(tags.SearchTagQuery{
		Name:         req.Name,
		FolderId:     req.FolderID,
		Limit:        &limit,
		Offset:       &offset,
		FoldersScope: &userResult.User.FolderIds,
	})
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	response := &SearchTagsResponse{
		Tags:   ConvertToTagsResponse(searchResult.Tags),
		Limit:  searchResult.Limit,
		Offset: searchResult.Offset,
		Total:  searchResult.Total,
	}

	return response, nil
}
