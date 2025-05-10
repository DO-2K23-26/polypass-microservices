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
	User types.User `json:"user" validate:"required"`
}

type UpdateUserRequest struct {
	ID        string `json:"id" validate:"required"`
	NewFolder string `json:"new_folder" validate:"required"`
}

type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}

type AddFolderAccessRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	FolderID string `json:"folder_id" validate:"required"`
}

type RemoveFolderAccessRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	FolderID string `json:"folder_id" validate:"required"`
}

type GetFoldersRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type GetFoldersResponse struct {
	Folders []types.Folder `json:"folders"`
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

func toAddFolderAccessResponse(result *user.AddFolderAccessResult) *UserResponse {
	if result == nil {
		return nil
	}

	return &UserResponse{
		User: result.User,
	}
}

func toRemoveFolderAccessResponse(result *user.RemoveFolderAccessResult) *UserResponse {
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
		User: req.User,
	}
}

func toUpdateUserQuery(req *UpdateUserRequest) user.UpdateUserQuery {
	return user.UpdateUserQuery{
		ID:        req.ID,
		NewFolder: req.NewFolder,
	}
}

func toDeleteUserQuery(req *DeleteUserRequest) user.DeleteUserQuery {
	return user.DeleteUserQuery{
		ID: req.ID,
	}
}

func toAddFolderAccessQuery(req *AddFolderAccessRequest) user.AddFolderAccessQuery {
	return user.AddFolderAccessQuery{
		UserID:   req.UserID,
		FolderID: req.FolderID,
	}
}

func toDeleteFolderAccessQuery(req *RemoveFolderAccessRequest) user.RemoveFolderAccessQuery {
	return user.RemoveFolderAccessQuery{
		UserID:   req.UserID,
		FolderID: req.FolderID,
	}
}
