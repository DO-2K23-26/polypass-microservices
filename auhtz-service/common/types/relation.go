package types

import "fmt"

type Relation string

var (
	// Used for folder
	Viewer  Relation = "viewer"
	Admin Relation = "admin"
	Parent Relation = "parent"
	
	// Used for credentials and tag
	FolderRelation Relation = "folder"
)


func ParseRelation(relationStr string) (Relation, error) {
	switch relationStr {
	case string(Viewer):
		return Viewer, nil
	case string(Admin):
		return Admin, nil
	case string(Parent):
		return Parent, nil
	case string(FolderRelation):
		return FolderRelation, nil
	default:
		return "", fmt.Errorf("invalid relation type: %s", relationStr)
	}
}