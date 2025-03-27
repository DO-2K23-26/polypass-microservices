package services

import (
	"sharing/dto"
	"sharing/repositories"
)

type SecretService interface {
	CreateSecret(secret dto.PostSecretRequest) (dto.PostSecretResponse, error)
	GetSecret(id string) (*dto.GetSecretResponse, error)
}

type secretService struct{
	repository repositories.SecretRepository
}

func NewSecretService() SecretService {
	return &secretService{
		repository: repositories.NewSecretRepository(),
	}
}

func (s *secretService) CreateSecret(secret dto.PostSecretRequest) (dto.PostSecretResponse, error) {

}