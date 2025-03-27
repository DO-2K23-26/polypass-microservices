package repositories

import (
	"sharing/models"
	"github.com/jmoiron/sqlx"
)

type SecretRepository interface {
	CreateSecret(secret *models.Secret) (models.Secret, error)
	GetSecret(id string) (*models.Secret, error)
}

func NewSecretRepository() SecretRepository {
	db := sqlx.MustConnect("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	return &secretRepository{db: db}
}

type secretRepository struct {
	db *sqlx.DB
}

func (r *secretRepository) CreateSecret(secret *models.Secret) (models.Secret, error) {
	query := `INSERT INTO secrets (content, expiration, is_encrypted, is_one_time_use, user_id) VALUES (:content, :expiration, :is_encrypted, :is_one_time_use, :user_id)`
	_, err := r.db.NamedExec(query, secret)
	if err != nil {
		return models.Secret{}, err
	}
	return *secret, nil
}

func (r *secretRepository) GetSecret(id string) (*models.Secret, error) {
	secret := models.Secret{}
	err := r.db.Get(&secret, "SELECT * FROM secrets WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}





