package graph

import (
	"github.com/DO-2K23-26/polypass-microservices/gateway/core"

	model1 "github.com/DO-2K23-26/polypass-microservices/gateway/graph/model"
)

func FolderToModel(folder *core.Folder) *model1.Folder {
	return &model1.Folder{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
		Icon:        folder.Icon,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}
}
