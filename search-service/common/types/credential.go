package types

type Credential struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	FolderId  string `json:"folder_id"`
	Tags      []Tag `json:"tags"`
	Folder    *Folder `json:"folder"`
}
