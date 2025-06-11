package graph

import "github.com/DO-2K23-26/polypass-microservices/gateway/core"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UsersService         core.UsersService
	OrganizationsService core.OrganizationService
}
