package app

import (
	"bytes"
	"time"
	"github.com/google/uuid"

	avroGeneratedSchema "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/generated"
	"github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/schemautils"
	"github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"gorm.io/gorm"
)

type TagService struct {
	publisher EventPublisher
	encoder   *schemautils.AvroEncoder
	database  *gorm.DB
}

func NewTagService(publisher EventPublisher, encoder *schemautils.AvroEncoder, database *gorm.DB) *TagService {
	return &TagService{publisher: publisher, encoder: encoder, database: database}
}

func (s *TagService) CreateTag(tag organization.CreateTagRequest) error {
	newTag := organization.Tag{
		Id:         uuid.New().String(),
		Name:       tag.Name,
		Color:      tag.Color,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		FolderID:   tag.FolderID,
		CreatedBy:  tag.CreatedBy,
	}

	data := avroGeneratedSchema.TagEvent{
		Id:         newTag.Id,
		Name:       newTag.Name,
		Color:      newTag.Color,
		Created_at:  newTag.CreatedAt.String(),
		Updated_at:  newTag.UpdatedAt.String(),
		Folder_id:   newTag.FolderID,
		Created_by:  newTag.CreatedBy,
	}

	var buf bytes.Buffer
	err := data.Serialize(&buf)
	if err != nil {
		return err
	}

	res := s.database.Create(&newTag)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return s.publisher.Publish("Tag-Create", buf.Bytes())
}

func (s *TagService) UpdateTag(tag organization.UpdateTagRequest) error {

	data := avroGeneratedSchema.TagEvent{
		Id:         tag.Id,
		Name:       tag.Name,
		Color:      tag.Color,
		Updated_at: time.Now().String(),
		Folder_id:  tag.FolderID,
	}

	var buf bytes.Buffer
	err := data.Serialize(&buf)
	if err != nil {
		return err
	}

	res := s.database.Model(&organization.Tag{}).
		Where("id = ?", tag.Id).
		Updates(organization.Tag{
			Name:       tag.Name,
			Color:      tag.Color,
			UpdatedAt:  time.Now(),
			FolderID:   tag.FolderID,
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return s.publisher.Publish("Tag-Update", buf.Bytes())
}

func (s *TagService) DeleteTag(id string) error {

	data := avroGeneratedSchema.TagEvent{
		Id:         id,
		Updated_at: time.Now().String(),
	}

	var buf bytes.Buffer
	err := data.Serialize(&buf)
	if err != nil {
		return err
	}

	res := s.database.Delete(&organization.Tag{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return s.publisher.Publish("tag-delete", buf.Bytes())
}

// Get all tags from the database.
func (s *TagService) ListTags(req organization.GetTagRequest) ([]organization.Tag, error) {
	var tags []organization.Tag

	if err := s.database.Model(&organization.Tag{}).Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Order("created_at asc").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// Get tag by its ID.
func (s *TagService) GetTag(id string) (*organization.Tag, error) {
	var tag organization.Tag
	res := s.database.Find(&tag, "id = ?", id)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &tag, nil
}
