package user

type IUserRepository interface {
	Get(query GetUserQuery) (*GetUserResult, error)
	Create(query CreateUserQuery) (*CreateUserResult, error)
	Delete(query DeleteUserQuery) (*DeleteUserResult, error)
	//By definition upsert action update the document if it exists, otherwise create a new one.
	//In this case, it will update a folder of the user
	AddFolderAccess(query AddFolderAccessQuery) (*AddFolderAccessResult, error)
	RemoveFolderAccess(query RemoveFolderAccessQuery) (*RemoveFolderAccessResult, error)
}
