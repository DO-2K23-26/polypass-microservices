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
	User types.User `json:"user"`
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