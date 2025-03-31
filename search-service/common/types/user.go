package types

type User struct {
	ID    string `json:"id"`
	Folders []Folder `json:"folders"`
}
