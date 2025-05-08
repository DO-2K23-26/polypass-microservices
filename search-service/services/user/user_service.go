package user

import (
	"errors"

	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"
)

type UserService struct {
	userRepository user.IUserRepository
}

func NewUserService(userRepository user.IUserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// GetUser retrieves a user by ID
func (s *UserService) Get(req *GetUserRequest) (*UserResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if req.ID == "" {
		return nil, errors.New("user ID is required")
	}

	query := toGetUserQuery(req)
	result, err := s.userRepository.Get(query)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("user not found")
	}

	return toUserResponse(result), nil
}

// CreateUser creates a new user
func (s *UserService) Create(req *CreateUserRequest) (*UserResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	query := toCreateUserQuery(req)
	result, err := s.userRepository.Create(query)
	if err != nil {
		return nil, err
	}

	return toCreateUserResponse(result), nil
}

// DeleteUser deletes a user by ID
func (s *UserService) Delete(req *DeleteUserRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.ID == "" {
		return errors.New("user ID is required")
	}

	query := toDeleteUserQuery(req)
	_, err := s.userRepository.Delete(query)
	if err != nil {
		return err
	}
	return nil
}

// AddFolderAccess adds or updates a user's access to a folder
func (s *UserService) GetFolders(req GetFoldersRequest) (*GetFoldersResponse, error) {

	if req.UserID == "" {
		return nil, errors.New("user ID is required")
	}

	result, err := s.userRepository.GetFolders(user.GetFoldersQuery{UserID: req.UserID})
	if err != nil {
		return nil, err
	}

	return &GetFoldersResponse{
		Folders: result.Folders,
	}, nil
}

// AddFolderAccess adds or updates a user's access to a folder
func (s *UserService) AddFolderAccess(req *AddFolderAccessRequest) (*UserResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if req.UserID == "" {
		return nil, errors.New("user ID is required")
	}

	if req.FolderID == "" {
		return nil, errors.New("folder ID is required")
	}

	query := toAddFolderAccessQuery(req)
	result, err := s.userRepository.AddFolderAccess(query)
	if err != nil {
		return nil, err
	}

	return toAddFolderAccessResponse(result), nil
}

// DeleteFolderAccess removes a user's access to a folder
func (s *UserService) RemoveFolderAccess(req *RemoveFolderAccessRequest) (*UserResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if req.UserID == "" {
		return nil, errors.New("user ID is required")
	}

	if req.FolderID == "" {
		return nil, errors.New("folder ID is required")
	}

	query := toDeleteFolderAccessQuery(req)
	result, err := s.userRepository.RemoveFolderAccess(query)
	if err != nil {
		return nil, err
	}

	return toRemoveFolderAccessResponse(result), nil
}
