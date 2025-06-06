package models

type Secret struct {
    Id           string            `json:"id" bson:"_id,omitempty" validate:"required,uuid"`
    Content      map[string]string `json:"content" bson:"content" validate:"required"`
    CreatedAt    int64             `json:"created_at" bson:"created_at" validate:"required,min=0"`
    Expiration   int64             `json:"expiration" bson:"expiration" validate:"required,min=0"`
    IsEncrypted  bool              `json:"is_encrypted" bson:"is_encrypted" validate:"required"`
    User         string            `json:"user_id" bson:"user_id"`
    IsOneTimeUse bool              `json:"is_one_time_use" bson:"is_one_time_use" validate:"required"`
    Name         string            `json:"name" bson:"name"`
}

type HistorySecret struct {
    Id           string            `json:"id" bson:"_id,omitempty" validate:"required,uuid"`
    CreatedAt    int64             `json:"created_at" bson:"created_at" validate:"required,min=0"`
    Expiration   int64             `json:"expiration" bson:"expiration" validate:"required,min=0"`
    ContentSize  int               `json:"content_size" bson:"content_size" validate:"required,min=0"`
    IsOneTimeUse bool              `json:"is_one_time_use" bson:"is_one_time_use" validate:"required"`
    Name         string            `json:"name" bson:"name"`
}