package tags

type ITagRepository interface {
	Get(query GetTagQuery) (*GetTagResult, error)
	Create(query CreateTagQuery) (*CreateTagResult, error)
	Update(query UpdateTagQuery) (*UpdateTagResult, error)
	Delete(query DeleteTagQuery) error
	SearchTags(query SearchTagQuery) (*SearchTagResult, error)
}
