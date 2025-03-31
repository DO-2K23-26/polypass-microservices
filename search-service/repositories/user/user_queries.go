package user

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

type GetUserQuery struct {
	ID string `json:"id"`
}

type GetUserResult struct {
	User types.User `json:"user"`
}

type CreateUserQuery struct {
	User types.User `json:"user"`
}

type CreateUserResult struct {
	User types.User `json:"user"`
}

type UpdateUserQuery struct {
	ID string `json:"id"`
	NewFolder string `json:"new_folder"` // This folder can be either the update of a folder or the user being added to a folder
}

type UpdateUserResult struct {
	User types.User `json:"user"`
}

type DeleteUserQuery struct {
	ID string `json:"id"`
}

type DeleteUserResult struct {
	ID string `json:"id"`
}

type UpsertFolderAccessQuery struct {
	UserID string `json:"user_id"`
	FolderID string `json:"folder_id"`
}

type UpsertFolderAccessResult struct {
	User types.User `json:"user"`
}

type DeleteFolderAccessQuery struct {
	UserID string `json:"user_id"`
	FolderID string `json:"folder_id"`
}

type DeleteFolderAccessResult struct {
	User types.User `json:"user"`
}