package models

type Secret struct {
	Id           string            `json:"id" validate:"required,uuid"`
	Content      map[string]string `json:"content" validate:"required"`
	Expiration   int64             `json:"expiration" validate:"required,min=0"`
	IsEncrypted  bool              `json:"is_encrypted" validate:"required"`
	IsOneTimeUse bool              `json:"is_one_time_use" validate:"required"`
}