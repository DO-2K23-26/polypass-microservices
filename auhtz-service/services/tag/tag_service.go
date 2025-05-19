package tag

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
)

type TagService struct {
	spiceDBAdapter *infrastructure.SpiceDBAdapter
}

func NewTagService(spiceDBAdapter *infrastructure.SpiceDBAdapter) *TagService {
	return &TagService{
		spiceDBAdapter: spiceDBAdapter,
	}
}

func (s *TagService) Create(ctx context.Context, tagID string, folderID string) error {
	// Create a relationship between the tag and the folder
	err := s.spiceDBAdapter.CreateRelationship(ctx, types.Folder, folderID, types.FolderRelation, types.Tag, tagID)
	if err != nil {
		log.Printf("failed to create tag relationship: %v", err)
		return err
	}

	log.Printf("successfully created tag relationship for tag %s and folder %s", tagID, folderID)
	return nil
}

func (s *TagService) Update(ctx context.Context, tagID string, folderID string) error {
	// Update the relationship between the tag and the folder
	err := s.spiceDBAdapter.UpdateRelationship(ctx, types.Folder, folderID, types.FolderRelation, types.Tag, tagID)
	if err != nil {
		log.Printf("failed to update tag relationship: %v", err)
		return err
	}

	log.Printf("successfully updated tag relationship for tag %s and folder %s", tagID, folderID)
	return nil
}

func (s *TagService) Delete(ctx context.Context, tagID string) error {
	// Delete the relationship for the tag
	err := s.spiceDBAdapter.Delete(ctx, types.Tag, tagID)
	if err != nil {
		log.Printf("failed to delete tag relationship: %v", err)
		return err
	}

	log.Printf("successfully deleted tag relationship for tag %s", tagID)
	return nil
}
