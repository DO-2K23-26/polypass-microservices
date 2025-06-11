package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// StringSlice est un type personnalisé pour gérer les tableaux de chaînes dans GORM
type StringSlice []string

// Scan implémente l'interface sql.Scanner pour le type StringSlice
func (ss *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*ss = StringSlice{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, ss)
	case string:
		return json.Unmarshal([]byte(v), ss)
	default:
		*ss = StringSlice{}
		return nil
	}
}

// Value implémente l'interface driver.Valuer pour le type StringSlice
func (ss StringSlice) Value() (driver.Value, error) {
	if ss == nil {
		return "[]", nil
	}
	return json.Marshal(ss)
}

type FolderEvent struct {
	ID          string   `avro:"id"`
	Name        string   `avro:"name"`
	Description string   `avro:"description"`
	Icon        string   `avro:"icon"`
	CreatedAt   string   `avro:"created_at"`
	UpdatedAt   string   `avro:"updated_at"`
	ParentID    string   `avro:"parent_id"`
	Members     []string `avro:"members"`
	CreatedBy   string   `avro:"created_by"`
}

// Folder représente la structure de dossier dans la base de données
type Folder struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ParentID    string `json:"parent_id"`
	Members     string `json:"members" gorm:"type:jsonb"` // Stocké comme une chaîne JSON
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ToFolderEvent convertit un Folder en FolderEvent
func (f *Folder) ToFolderEvent() *FolderEvent {
	var members []string
	if f.Members != "" {
		json.Unmarshal([]byte(f.Members), &members)
	}
	return &FolderEvent{
		ID:          f.ID,
		Name:        f.Name,
		Description: f.Description,
		Icon:        f.Icon,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
		ParentID:    f.ParentID,
		Members:     members,
		CreatedBy:   f.CreatedBy,
	}
}

// FromFolderEvent met à jour un Folder à partir d'un FolderEvent
func (f *Folder) FromFolderEvent(event *FolderEvent) {
	f.ID = event.ID
	f.Name = event.Name
	f.Description = event.Description
	f.Icon = event.Icon
	f.CreatedAt = event.CreatedAt
	f.UpdatedAt = event.UpdatedAt
	f.ParentID = event.ParentID
	f.CreatedBy = event.CreatedBy

	// Convertir les membres en JSON
	if event.Members != nil {
		data, _ := json.Marshal(event.Members)
		f.Members = string(data)
	} else {
		f.Members = "[]"
	}
}

// GetMembers retourne les membres sous forme de slice de string
func (f *Folder) GetMembers() []string {
	var members []string
	if f.Members != "" {
		json.Unmarshal([]byte(f.Members), &members)
	}
	return members
}

// SetMembers définit les membres à partir d'un slice de string
func (f *Folder) SetMembers(members []string) {
	if members == nil {
		f.Members = "[]"
		return
	}
	data, _ := json.Marshal(members)
	f.Members = string(data)
}

var EsFolder = map[string]types.Property{
	"id":          types.NewKeywordProperty(),
	"name":        types.NewSearchAsYouTypeProperty(),
	"description": types.NewTextProperty(),
	"icon":        types.NewKeywordProperty(),
	"parent_id":   types.NewKeywordProperty(),
	"members":     types.NewKeywordProperty(),
	"created_by":  types.NewKeywordProperty(),
	"created_at":  types.NewDateProperty(),
	"updated_at":  types.NewDateProperty(),
}

// The parent id is stringified
type FolderSql struct {
	ID       string         `json:"id"  gorm:"primaryKey"`
	Name     string         `json:"name"`
	ParentID sql.NullString `json:"parent_id"`
}

var FolderIndex = "folders"
