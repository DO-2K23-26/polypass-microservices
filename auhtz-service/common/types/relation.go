package types

type Relation string

var (
	// Used for folder
	Viewer  Relation = "viewer"
	Admin Relation = "admin"
	Parent Relation = "parent"
	
	// Used for credentials and tag
	FolderRelation Relation = "folder"
)
