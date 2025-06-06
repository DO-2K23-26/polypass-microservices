package services

import (
	"fmt"
	"sharing/dto"
	"sharing/models"
	"sharing/repositories"
	"time"
)

type SecretService interface {
	CreateSecret(secret dto.PostSecretRequest, userId string) (dto.PostSecretResponse, error)
	GetSecret(id string) (*dto.GetSecretResponse, error)
	GetHistory(userId string) (*[]models.HistorySecret, error)
}

type secretService struct{
	repository repositories.SecretRepository
}

func NewSecretService() SecretService {
	return &secretService{
		repository: repositories.NewSecretRepository(),
	}
}

func (s *secretService) CreateSecret(secret dto.PostSecretRequest, userId string) (dto.PostSecretResponse, error) {
	secretEntity := models.Secret{
		Content:      secret.Content,
		Name:         secret.Name,
		Expiration:   secret.Expiration,
		CreatedAt:    time.Now().Unix(),
		IsEncrypted:  secret.IsEncrypted,
		IsOneTimeUse: secret.IsOneTimeUse,
		User:         userId,
	}

	createdSecret, err := s.repository.CreateSecret(&secretEntity)

	if err != nil {
		return dto.PostSecretResponse{}, err
	}

	createdSecretResponse := dto.PostSecretResponse{
		Id:           createdSecret.Id,
		CreatedAt: 		time.Now().Unix(),
	}

	return createdSecretResponse, nil
}

func (s *secretService) GetSecret(id string) (*dto.GetSecretResponse, error) {
	secret, err := s.repository.GetSecret(id)

	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, fmt.Errorf("Secret not found")
	}

	secretResponse := dto.GetSecretResponse{
		Id:          secret.Id,
		IsEncrypted: secret.IsEncrypted,
		Content:     secret.Content,
	}

	return &secretResponse, nil
}

func (s *secretService) GetHistory(userId string) (*[]models.HistorySecret, error) {
	history, err := s.repository.GetHistory(userId)

	if err != nil {
		return nil, err
	}

	if history == nil {
		return nil, fmt.Errorf("No history found")
	}

	return &history, nil
}