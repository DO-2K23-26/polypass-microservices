package users

type UserAPI interface {
	Setup() error
	Shutdown() error
	GetUser(id string) (*GetUserResponse, error)
}

type GetUserResponse struct {
	ID        string `json:"Id"`
	Username  string `json:"Username"`
	Email     string `json:"Email"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
}

type usersAPI struct {

}

// func NewUsersAPI(config UsersConfig) UserAPI {
// 	return &usersAPI{}
// }
