package types

type User struct {
	ID    string `json:"id"`
	FolderIds []string `json:"folder_ids"`
}
