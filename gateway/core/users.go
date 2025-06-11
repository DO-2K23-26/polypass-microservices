package core


import "github.com/DO-2K23-26/polypass-microservices/gateway/infrastructure/users"



type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}	

type UsersService interface {
	GetUser(id string) (User, error)
}

type usersService struct {
	usersAPI users.UserAPI
}

func NewUsersService(usersAPI users.UserAPI) UsersService {
	return &usersService{
		usersAPI: usersAPI,
	}
}

func (s *usersService) GetUser(id string) (User, error) {
	response, err := s.usersAPI.GetUser(id)

	if err != nil {
		return User{}, err
	}

	return FromGetUserResponse(response), nil
}

func FromGetUserResponse(response *users.GetUserResponse) User {
	return User{
		ID:        response.ID,
		Username:  response.Username,
		Email:     response.Email,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	}
}
