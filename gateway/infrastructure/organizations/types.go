package organizations

type CreateFolderRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ParentID    string `json:"parent_id"`
	CreatedBy   string `json:"created_by"`
}

type CreateFolderResponse struct {
	ID          string   `json:"Id"`
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Icon        string   `json:"Icon"`
	CreatedAt   string   `json:"CreatedAt"`
	UpdatedAt   string   `json:"UpdatedAt"`
	ParentID    *string  `json:"ParentID"`
	CreatedBy   string   `json:"CreatedBy"`
	Members     []string `json:"members"`
}

type GetFolderResponse struct {
	ID          string   `json:"Id"`
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Icon        string   `json:"Icon"`
	CreatedAt   string   `json:"CreatedAt"`
	UpdatedAt   string   `json:"UpdatedAt"`
	ParentID    *string  `json:"ParentID"`
	CreatedBy   string   `json:"CreatedBy"`
	Members     []string `json:"members"`
}



type GetFoldersResponse struct {
	Folders []GetFolderResponse `json:"folders"`
	Total   int                 `json:"total"`
	Page    int                 `json:"page"`
	Limit   int                 `json:"limit"`
}

type UpdateFolderRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
