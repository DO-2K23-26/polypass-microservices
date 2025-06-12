package core

type Folder struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Icon string `json:"icon"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ParentID *string `json:"parent_id"`
	CreatedBy string `json:"created_by"`
	Members []string `json:"members"`
}
