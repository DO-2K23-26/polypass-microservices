package user

type UserRepository interface {
	Get(query GetUserQuery) (*GetUserResult, error)
	Create(query CreateUserQuery) (*CreateUserResult, error)
	Update(query UpdateUserQuery) (*UpdateUserResult, error)
	Delete(query DeleteUserQuery) (*DeleteUserResult, error)
}
