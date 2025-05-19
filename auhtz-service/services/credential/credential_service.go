package credential

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
)

type CredentialService struct {
	spiceDBAdapter *infrastructure.SpiceDBAdapter
}

func NewCredentialService(spiceDBAdapter *infrastructure.SpiceDBAdapter) *CredentialService {
	return &CredentialService{
		spiceDBAdapter: spiceDBAdapter,
	}
}

func (s *CredentialService) Create(ctx context.Context, credentialID string, folderID string) error {
	// Create a relationship in SpiceDB
	err := s.spiceDBAdapter.CreateRelationship(ctx, types.Folder, folderID, types.FolderRelation, types.Credential, credentialID)
	if err != nil {
		log.Printf("failed to create credential relationship: %v", err)
		return err
	}

	log.Printf("successfully created credential relationship for credential %s in folder %s", credentialID, folderID)
	return nil
}

func (s *CredentialService) Update(ctx context.Context, credentialID string, folderID string) error {
	// Update the relationship in SpiceDB
	err := s.spiceDBAdapter.UpdateRelationship(ctx, types.Folder, folderID, types.FolderRelation, types.Credential, credentialID)
	if err != nil {
		log.Printf("failed to update credential relationship: %v", err)
		return err
	}

	log.Printf("successfully updated credential relationship for credential %s in folder %s", credentialID, folderID)
	return nil
}

func (s *CredentialService) Delete(ctx context.Context, credentialID string) error {
	// Delete the relationship in SpiceDB
	err := s.spiceDBAdapter.Delete(ctx, types.Credential, credentialID)
	if err != nil {
		log.Printf("failed to delete credential relationship: %v", err)
		return err
	}

	log.Printf("successfully deleted credential relationship for credential %s", credentialID)
	return nil
}
