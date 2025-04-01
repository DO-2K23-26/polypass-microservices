package models

type Secret struct {
    Id           string            `json:"id" bson:"_id,omitempty" validate:"required,uuid"`
    Content      map[string]string `json:"content" bson:"content" validate:"required"`
    Expiration   int64             `json:"expiration" bson:"expiration" validate:"required,min=0"`
    IsEncrypted  bool              `json:"is_encrypted" bson:"is_encrypted" validate:"required"`
    IsOneTimeUse bool              `json:"is_one_time_use" bson:"is_one_time_use" validate:"required"`
}