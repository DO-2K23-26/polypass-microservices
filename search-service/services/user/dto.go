package user

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"
)

// Request DTOs
type GetUserRequest struct {
	ID string `json:"id" validate:"required"`
}

type CreateUserRequest struct {
	ID string `json:"id" validate:"required"`
}

type UpdateUserRequest struct {
	ID        string `json:"id" validate:"required"`
	NewFolder string `json:"newFolder" validate:"required"`
}

type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// Response DTOs
type UserResponse struct {
	User types.User `json:"user"`
}

// Conversion methods for repository to service layer
func toUserResponse(result *user.GetUserResult) *UserResponse {
	if result == nil {
		return nil
	}
	
	return &UserResponse{
		User: result.User,
	}
}

func toCreateUserResponse(result *user.CreateUserResult) *UserResponse {
	if result == nil {
		return nil
	}
	
	return &UserResponse{
		User: result.User,
	}
}

func toUpdateUserResponse(result *user.UpdateUserResult) *UserResponse {
	if result == nil {
		return nil
	}
	
	return &UserResponse{
		User: result.User,
	}
}

// Conversion methods for service to repository layer
func toGetUserQuery(req *GetUserRequest) user.GetUserQuery {
	return user.GetUserQuery{
		ID: req.ID,
		
	}
}

func toCreateUserQuery(req *CreateUserRequest) user.CreateUserQuery {
	return user.CreateUserQuery{
		User: types.User{
			ID: req.ID,
		},
	}
}

// func toUpdateUserQuery(req *UpdateUserRequest) user.UpdateUserQuery {
// 	return user.UpdateUserQuery{
// 		User: types.User{
// 			ID:     req.ID,
			
// 		},
// 	}
// }

func toDeleteUserQuery(req *DeleteUserRequest) user.DeleteUserQuery {
	return user.DeleteUserQuery{
		ID: req.ID,
	}
}