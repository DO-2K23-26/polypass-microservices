package folder

type IFolderRepository interface {
	CreateFolder(query CreateFolderQuery) (*CreateFolderResult, error)
	GetFolder(query GetFolderQuery) (*GetFolderResult, error)
	UpdateFolder(query UpdateFolderQuery) (*UpdateFolderResult, error)
	DeleteFolder(query DeleteFolderQuery) error
	GetFolderHierarchy(query GetFolderHierarchyQuery) (*GetFolderHierarchyResult, error)
	SearchFolder(query SearchFolderQuery) (*SearchFolderResult, error)
}
