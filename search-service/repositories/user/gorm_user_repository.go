package user

import (
	"errors"

	"github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	DB *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		DB: db,
	}
}

func (r *GormUserRepository) Get(query GetUserQuery) (*GetUserResult, error) {
	var user types.User
	if err := r.DB.Where("id = ?", query.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &GetUserResult{
		User: user,
	}, nil
}

func (r *GormUserRepository) Create(query CreateUserQuery) (*CreateUserResult, error) {
	if err := r.DB.Create(&query.User).Error; err != nil {
		return nil, err
	}

	return &CreateUserResult{
		User: query.User,
	}, nil
}

func (r *GormUserRepository) Delete(query DeleteUserQuery) (*DeleteUserResult, error) {
	if err := r.DB.Delete(&types.User{}, "id = ?", query.ID).Error; err != nil {
		return nil, err
	}

	return &DeleteUserResult{
		ID: query.ID,
	}, nil
}

func (r *GormUserRepository) AddFolderAccess(query AddFolderAccessQuery) (*AddFolderAccessResult, error) {
	var user types.User
	if err := r.DB.Where("id = ?", query.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check if folder exists in the database
	var folder types.Folder
	if err := r.DB.Where("id = ?", query.FolderID).First(&folder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("folder not found")
		}
		return nil, err
	}

	// Check if folder already exists in user's folders
	for _, userFolder := range user.Folders {
		if userFolder.ID == query.FolderID {
			return &AddFolderAccessResult{User: user}, nil
		}
	}

	// Add the folder to the user's folders
	user.Folders = append(user.Folders, folder)
	if err := r.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &AddFolderAccessResult{
		User: user,
	}, nil
}

func (r *GormUserRepository) RemoveFolderAccess(query RemoveFolderAccessQuery) (*RemoveFolderAccessResult, error) {
	var user types.User
	if err := r.DB.Where("id = ?", query.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Remove the folder from user's folders
	updatedFolders := make([]types.Folder, 0)
	for _, folder := range user.Folders {
		if folder.ID != query.FolderID {
			updatedFolders = append(updatedFolders, folder)
		}
	}

	// Only update if there was a change
	if len(updatedFolders) != len(user.Folders) {
		user.Folders = updatedFolders
		if err := r.DB.Save(&user).Error; err != nil {
			return nil, err
		}
	}

	return &RemoveFolderAccessResult{
		User: user,
	}, nil
}
