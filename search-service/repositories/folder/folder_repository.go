package folder

type IFolderRepository interface {
	Create(query CreateFolderQuery) (*CreateFolderResult, error)
	Get(query GetFolderQuery) (*GetFolderResult, error)
	Update(query UpdateFolderQuery) (*UpdateFolderResult, error)
	Delete(query DeleteFolderQuery) error
	GetHierarchy(query GetFolderHierarchyQuery) (*GetFolderHierarchyResult, error)
	Search(query SearchFolderQuery) (*SearchFolderResult, error)
}
