package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	avroGeneratedSchema "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/generated"
	"github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/schemautils"
	organization "github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"gorm.io/gorm"
)

// FolderCredentialService links folders with credentials using the credential service.
type FolderCredentialService struct {
	db        *gorm.DB
	host      string
	client    *http.Client
	publisher EventPublisher
	encoder   *schemautils.AvroEncoder
}

// NewFolderCredentialService creates a new FolderCredentialService.
func NewFolderCredentialService(db *gorm.DB, publisher EventPublisher, encoder *schemautils.AvroEncoder) *FolderCredentialService {
	host := os.Getenv("CREDENTIAL_SERVICE_HOST")
	if host == "" {
		host = "http://localhost:8080"
		log.Println("CREDENTIAL_SERVICE_HOST is not set, using default value (http://localhost:8080)")
	}
	return &FolderCredentialService{db: db, host: host, client: &http.Client{}, publisher: publisher, encoder: encoder}
}

// List returns paginated credentials for a folder.
func (s *FolderCredentialService) List(folderID, credType string, page, limit int) (*organization.CredentialList, error) {
	var relations []organization.FolderCredential
	if err := s.db.Where("id_folder = ?", folderID).Find(&relations).Error; err != nil {
		return nil, err
	}

	start := (page - 1) * limit
	if start > len(relations) {
		return &organization.CredentialList{Credentials: []map[string]interface{}{}, Page: page, Limit: limit}, nil
	}
	end := start + limit
	if end > len(relations) {
		end = len(relations)
	}

	creds := make([]map[string]interface{}, 0, end-start)
	for _, rel := range relations[start:end] {
		url := fmt.Sprintf("%s/credentials/%s/%s", s.host, credType, rel.IdCredential)
		resp, err := s.client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("credential service returned %d", resp.StatusCode)
		}
		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, err
		}
		creds = append(creds, data)
	}

	return &organization.CredentialList{Credentials: creds, Page: page, Limit: limit}, nil
}

// Create creates a credential via the credential service and stores the link.
func (s *FolderCredentialService) Create(folderID, credType string, body []byte) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/credentials/%s", s.host, credType)
	resp, err := s.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("credential service returned %d: %s", resp.StatusCode, string(b))
	}
	var credential map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&credential); err != nil {
		return nil, err
	}
	id, ok := credential["id"].(string)
	if !ok {
		return nil, fmt.Errorf("credential id not found in response")
	}
	rel := organization.FolderCredential{IdFolder: folderID, IdCredential: id}
	if err := s.db.Create(&rel).Error; err != nil {
		return nil, err
	}

	name, _ := credential["name"].(string)
	event := avroGeneratedSchema.CredentialEvent{
		Credential_id:   id,
		Credential_name: name,
		Folder_id:       folderID,
	}
	var buf bytes.Buffer
	if err := event.Serialize(&buf); err != nil {
		return nil, err
	}
	if err := s.publisher.Publish("credential-creation", buf.Bytes()); err != nil {
		return nil, err
	}

	return credential, nil
}

// Update updates a credential via the credential service.
func (s *FolderCredentialService) Update(folderID, credType, credentialID string, body []byte) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/credentials/%s/%s", s.host, credType, credentialID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("credential service returned %d: %s", resp.StatusCode, string(b))
	}
	var credential map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&credential); err != nil {
		return nil, err
	}

	name, _ := credential["name"].(string)
	event := avroGeneratedSchema.CredentialEvent{
		Credential_id:   credentialID,
		Credential_name: name,
		Folder_id:       folderID,
	}
	var buf bytes.Buffer
	if err := event.Serialize(&buf); err != nil {
		return nil, err
	}
	if err := s.publisher.Publish("credential-update", buf.Bytes()); err != nil {
		return nil, err
	}

	return credential, nil
}

// Delete removes credentials via the credential service and unlinks them from the folder.
func (s *FolderCredentialService) Delete(folderID, credType string, ids []string) error {
	names := make(map[string]string, len(ids))
	for _, id := range ids {
		url := fmt.Sprintf("%s/credentials/%s/%s", s.host, credType, id)
		resp, err := s.client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var data map[string]interface{}
				if json.NewDecoder(resp.Body).Decode(&data) == nil {
					if n, ok := data["name"].(string); ok {
						names[id] = n
					}
				}
			} else {
				io.ReadAll(resp.Body)
			}
		}
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/credentials/%s", s.host, credType), nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	for _, id := range ids {
		q.Add("id", id)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("credential service returned %d: %s", resp.StatusCode, string(b))
	}

	if err := s.db.Where("id_folder = ? AND id_credential IN ?", folderID, ids).Delete(&organization.FolderCredential{}).Error; err != nil {
		return err
	}

	for _, id := range ids {
		event := avroGeneratedSchema.CredentialEvent{
			Credential_id:   id,
			Credential_name: names[id],
			Folder_id:       folderID,
		}
		var buf bytes.Buffer
		if err := event.Serialize(&buf); err != nil {
			return err
		}
		if err := s.publisher.Publish("credential-delete", buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}
