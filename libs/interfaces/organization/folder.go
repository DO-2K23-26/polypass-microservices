package organization

import (
	"time"

	"github.com/lib/pq"
)

type Folder struct {
	Id          string `gorm:"primaryKey;type:uuid"`
	Name        string `gorm:"not null"`
	Description *string
	Icon        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ParentID    *string
	CreatedBy   string
	Members     pq.StringArray `gorm:"type:text[]" json:"members"`
}

type CreateFolderRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	ParentID    *string `json:"parent_id"`
	CreatedBy   string  `json:"created_by" binding:"required"`
}
