package folder

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
)

type FolderService struct {
	spiceDBAdapter *infrastructure.SpiceDBAdapter
}

func NewFolderService(spiceDBAdapter *infrastructure.SpiceDBAdapter) *FolderService {
	return &FolderService{
		spiceDBAdapter: spiceDBAdapter,
	}
}

func (s *FolderService) Create(ctx context.Context, folderId string, parentId string) error {
	err := s.spiceDBAdapter.CreateRelationship(ctx, types.Folder, folderId, types.Parent, types.Folder, parentId)
	if err != nil {
		log.Printf("failed to create folder relationship: %v", err)
		return err
	}
	log.Printf("successfully created folder %s with parent %s", folderId, parentId)
	return nil
}

func (s *FolderService) Delete(ctx context.Context, folderId string) error {
	err := s.spiceDBAdapter.Delete(ctx, types.Folder, folderId)
	if err != nil {
		log.Printf("failed to delete folder relationships: %v", err)
		return err
	}
	log.Printf("successfully deleted folder %s", folderId)
	return nil
}

func (s *FolderService) UpdateParent(ctx context.Context, folderId string, newParentId string) error {
	err := s.spiceDBAdapter.ChangeParent(ctx, folderId, newParentId)
	if err != nil {
		log.Printf("failed to update folder parent: %v", err)
		return err
	}
	log.Printf("successfully updated folder %s to new parent %s", folderId, newParentId)
	return nil
}
