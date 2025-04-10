package folder

type FolderRepository interface {
	CreateFolder(query CreateFolderQuery) (*CreateFolderResult, error)
	GetFolder(query GetFolderQuery) (*GetFolderResult, error)
	UpdateFolder(query UpdateFolderQuery) (*UpdateFolderResult, error)
	DeleteFolder(query DeleteFolderQuery) error
	SearchFolder(query SearchFolderQuery) (*SearchFolderResult, error)
}
