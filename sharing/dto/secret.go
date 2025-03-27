package dto

type PostSecretRequest struct {
	Content       map[string]string `json:"content" validate:"required"`
	Expiration    int64             `json:"expiration" validate:"required,min=0"`
	IsEncrypted   bool              `json:"is_encrypted" validate:"required"`
	IsOneTimeUse  bool              `json:"is_one_time_use" validate:"required"`
}

type PostSecretResponse struct {
	Id string `json:"id"`
	CreatedAt int64 `json:"created_at"`
}

type GetSecretResponse struct {
	Id          string            `json:"id" validate:"required,uuid"`
	IsEncrypted bool              `json:"is_encrypted" validate:"required"`
	Content     map[string]string `json:"content" validate:"required"`
}