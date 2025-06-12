package organizations

import (
	"strconv"

	"resty.dev/v3"
)

type OrganizationAPI interface {
	Setup() error
	Shutdown() error
	CreateFolder(payload CreateFolderRequest) (*CreateFolderResponse, error)
	GetFolder(id string) (*GetFolderResponse, error)
	GetFolders(page int, limit int) (*GetFoldersResponse, error)
	UpdateFolder(id string, payload CreateFolderRequest) (*GetFolderResponse, error)
	DeleteFolder(id string) error
}

type organizationsAPI struct {
	client  *resty.Client
	baseURL string
}

func New(config Config) OrganizationAPI {
	c := resty.New()
	return &organizationsAPI{
		client:  c,
		baseURL: config.Endpoint,
	}
}

func (o *organizationsAPI) Setup() error {
	return nil
}

func (o *organizationsAPI) Shutdown() error {
	return o.client.Close()
}

func (o *organizationsAPI) CreateFolder(payload CreateFolderRequest) (*CreateFolderResponse, error) {
	res, err := o.client.R().
		SetBody(payload).
		SetResult(&CreateFolderResponse{}).
		Post(o.baseURL + "/folders")
	if err != nil {
		return nil, err
	}
	response := res.Result().(*CreateFolderResponse)
	return response, nil
}

func (o *organizationsAPI) GetFolder(id string) (*GetFolderResponse, error) {
	res, err := o.client.R().
		SetResult(&GetFolderResponse{}).
		Get(o.baseURL + "/folders/" + id)
	if err != nil {
		return nil, err
	}
	response := res.Result().(*GetFolderResponse)
	return response, nil
}

func (o *organizationsAPI) GetFolders(page int, limit int) (*GetFoldersResponse, error) {
	res, err := o.client.R().
		SetResult(&GetFoldersResponse{}).
		Get(o.baseURL + "/folders?page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}
	response := res.Result().(*GetFoldersResponse)
	return response, nil
}

func (o *organizationsAPI) UpdateFolder(id string, payload CreateFolderRequest) (*GetFolderResponse, error) {
	res, err := o.client.R().
		SetBody(payload).
		SetResult(&GetFolderResponse{}).
		Put(o.baseURL + "/folders/" + id)
	if err != nil {
		return nil, err
	}
	response := res.Result().(*GetFolderResponse)
	return response, nil
}

func (o *organizationsAPI) DeleteFolder(id string) error {
	_, err := o.client.R().
		Delete(o.baseURL + "/folders/" + id)
	if err != nil {
		return err
	}
	return nil
}
