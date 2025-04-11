package user

import "github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"

type ElasticUserRepository struct {
	EsClient infrastructure.ElasticAdapter
}

// Implements the UserRepository interface at ./user_repository.go.  