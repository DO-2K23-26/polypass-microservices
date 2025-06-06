package organization

import (
	"time"
	"github.com/lib/pq"
)

type Folder struct {
	ID          string         `gorm:"primaryKey;type:uuid"`
	Name        string         `gorm:"not null"`
	Description *string
	Icon        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ParentID    *string
	CreatedBy   string
	Members     pq.StringArray `gorm:"type:text[]" json:"members"`
}
