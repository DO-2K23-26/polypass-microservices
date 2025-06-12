package core

import "github.com/DO-2K23-26/polypass-microservices/gateway/infrastructure/organizations"

type OrganizationService interface {
	CreateFolder(request organizations.CreateFolderRequest) (Folder, error)
	GetFolder(id string) (Folder, error)
	GetFolders() ([]Folder, error)
	DeleteFolder(id string) error
	UpdateFolder(id string, request organizations.CreateFolderRequest) (Folder, error)
}

type organizationsService struct {
	organizationAPI organizations.OrganizationAPI
}

func FromCreateFolderResponse(response *organizations.CreateFolderResponse) Folder {
	return Folder{
		ID:          response.ID,
		Name:        response.Name,
		Description: response.Description,
		Icon:        response.Icon,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
		ParentID:    response.ParentID,
		CreatedBy:   response.CreatedBy,
		Members:     response.Members,
	}
}

func FromGetFolderResponse(response *organizations.GetFolderResponse) Folder {
	return Folder{
		ID:          response.ID,
		Name:        response.Name,
		Description: response.Description,
		Icon:        response.Icon,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
		ParentID:    response.ParentID,
		CreatedBy:   response.CreatedBy,
		Members:     response.Members,
	}
}

func FromGetFoldersResponse(response *organizations.GetFoldersResponse) []Folder {
	var folders []Folder

	for _, folder := range response.Folders {
		folders = append(folders, FromGetFolderResponse(&folder))
	}

	return folders
}

func NewOrganizationsService(organizationAPI organizations.OrganizationAPI) OrganizationService {
	return &organizationsService{
		organizationAPI: organizationAPI,
	}
}

func (s *organizationsService) CreateFolder(request organizations.CreateFolderRequest) (Folder, error) {
	response, err := s.organizationAPI.CreateFolder(request)

	if err != nil {
		return Folder{}, err
	}

	return FromCreateFolderResponse(response), nil
}

func (s *organizationsService) GetFolder(id string) (Folder, error) {
	response, err := s.organizationAPI.GetFolder(id)

	if err != nil {
		return Folder{}, err
	}

	return FromGetFolderResponse(response), nil
}

func (s *organizationsService) GetFolders() ([]Folder, error) {
	response, err := s.organizationAPI.GetFolders(1,100)

	if err != nil {
		return []Folder{}, err
	}

	return FromGetFoldersResponse(response), nil
}

func (s *organizationsService) DeleteFolder(id string) error {
	return s.organizationAPI.DeleteFolder(id)
}

func (s *organizationsService) UpdateFolder(id string, request organizations.CreateFolderRequest) (Folder, error) {
	response, err := s.organizationAPI.UpdateFolder(id, request)

	if err != nil {
		return Folder{}, err
	}

	return FromGetFolderResponse(response), nil
}


