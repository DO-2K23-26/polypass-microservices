package user

import (
	"errors"

	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"
)

type UserService struct {
	userRepository user.UserRepository
}

func NewUserService(userRepository user.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}



// GetUser retrieves a user by ID
func (s *UserService) GetUser(req *GetUserRequest) (*UserResponse, error) {
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
func (s *UserService) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
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

// // UpdateUser updates an existing user
// func (s *UserService) UpdateUser(req *UpdateUserRequest) (*UserResponse, error) {
// 	if req == nil {
// 		return nil, errors.New("request cannot be nil")
// 	}

// 	if req.ID == "" {
// 		return nil, errors.New("user ID is required")
// 	}

// 	query := toUpdateUserQuery(req)
// 	result, err := s.userRepository.Update(query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return toUpdateUserResponse(result), nil
// }

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(req *DeleteUserRequest) (string, error) {
	if req == nil {
		return "", errors.New("request cannot be nil")
	}

	if req.ID == "" {
		return "", errors.New("user ID is required")
	}

	query := toDeleteUserQuery(req)
	result, err := s.userRepository.Delete(query)
	if err != nil {
		return "", err
	}

	return result.ID, nil
}
