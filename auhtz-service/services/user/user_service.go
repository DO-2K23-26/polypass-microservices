package user

import (
	"context"
	"fmt"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
)

type UserService struct {
	spiceDBAdapter *infrastructure.SpiceDBAdapter
}

func NewUserService(spiceDBAdapter *infrastructure.SpiceDBAdapter) *UserService {
	return &UserService{
		spiceDBAdapter: spiceDBAdapter,
	}
}

func (s *UserService) AddUserToFolder(ctx context.Context, userId string, folderId string, relation string) error {
	parsedRelation, err := types.ParseRelation(relation)
	if err != nil {
		return err
	}
	if parsedRelation != types.Viewer && parsedRelation != types.Admin {
		return fmt.Errorf("invalid relation type: %s", relation)
	}
	return s.spiceDBAdapter.CreateRelationship(ctx, types.User, userId, parsedRelation, types.Folder, folderId)
}

func (s *UserService) RemoveUserFromFolder(ctx context.Context, userId string, folderId string) error {
	return s.spiceDBAdapter.DeleteRelationship(
		ctx,
		types.User,
		userId,
		types.Folder,
		folderId,
	)

}
