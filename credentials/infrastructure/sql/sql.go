package sql

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/DO-2K23-26/polypass-microservices/credentials/types"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/linkedin/goavro"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Sql interface {
	Setup() error
	Shutdown() error
	// pass an empty interface to get a credential from database

	GetPasswordCredentials(ids []string) ([]types.PasswordCredential, error)
	GetCardCredentials(ids []string) ([]types.CardCredential, error)
	GetSSHKeyCredentials(ids []string) ([]types.SSHKeyCredential, error)

	CreatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error)
	CreateCardCredential(credential types.CardCredential) (types.CardCredential, error)
	CreateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error)

	UpdatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error)
	UpdateCardCredential(credential types.CardCredential) (types.CardCredential, error)
	UpdateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error)

	DeletePasswordCredentials(ids []string) error
	DeleteCardCredentials(ids []string) error
	DeleteSSHKeyCredentials(ids []string) error
}

type sql struct {
	db         *sqlx.DB
	producer   *kafka.Producer
	consumer   *kafka.Consumer
	migrations string
	username   string
	password   string
	host       string
	port       int
	dbname     string
}

func NewSql(config Config,producer *kafka.Producer, consumer *kafka.Consumer) (Sql, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Dbname))
	if err != nil {
		return nil, err
	}
	return sql{
		db:         db,
		producer:   producer,
		consumer: consumer,
		migrations: config.Migrations,
		username:   config.Username,
		password:   config.Password,
		host:       config.Host,
		port:       config.Port,
		dbname:     config.Dbname,
	}, nil
}

func (m sql) Setup() error {
	migrations, err := migrate.New(fmt.Sprintf("file://%s", m.migrations), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", m.username, m.password, m.host, m.port, m.dbname))

	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil {
		return err
	}
	fmt.Println("Migrations up")

	return nil
}

func (m sql) Shutdown() error {
	return m.db.Close()
}

// Define a mapping from credentials types to schema file paths
var credentialsSchema = map[string]string{
	"PasswordCredential": "interfaces/credentials/password_credential.avsc",
	"CardCredential":     "interfaces/credentials/card_credential.avsc",
	"SSHKeyCredential":   "interfaces/credentials/ssh_credential.avsc",
	"CardAttribute":     "interfaces/credentials/card_attributes.avsc",
	"CreateCredentialsOpts": "interfaces/credentials/create_credentials_opts.avsc",
	"Credential": "interfaces/credentials/credential.avsc",
	"PasswordAttribute": "interfaces/credentials/password_attributes.avsc",
	"SSHKeyAttribute": "interfaces/credentials/ssh_attribute.avsc",
	"UserIdentifierAttributes": "interfaces/credentials/user_identifier_attribues.avsc",
}

// Load the Avro schema from a file
func loadSchema(filePath string) (string, error) {
	schemaBytes, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}
	return string(schemaBytes), nil
}

// Get the schema path based on the topic
func getSchemaPath(topic string) (string, error) {
	schemaPath, exists := credentialsSchema[topic]
	if !exists {
		return "", fmt.Errorf("no schema found for topic: %s", topic)
	}
	return schemaPath, nil
}

func (m *sql) produceMessage(topic string, key string, value []byte) error {
	// Get the schema path based on the topic
	schemaPath, err := getSchemaPath(topic)
	if err != nil {
		log.Printf("Failed to get schema path: %v", err)
		return err
	}

	// Load the appropriate Avro schema
	schema, err := loadSchema(schemaPath)
	if err != nil {
		log.Printf("Failed to load schema: %v", err)
		return err
	}

	// Create a new Avro codec
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		log.Printf("Failed to create codec: %v", err)
		return err
	}

	// Create a map to hold the message data
	messageMap := map[string]interface{}{
		"id":      key,
		"message": string(value),
	}

	// Serialize the message using Avro
	avroBinary, err := codec.BinaryFromNative(nil, messageMap)
	if err != nil {
		log.Printf("Failed to serialize data: %v", err)
		return err
	}

	// Produce the message to Kafka
	deliveryChan := make(chan kafka.Event)
	err = m.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          avroBinary,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	msg := e.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		return msg.TopicPartition.Error
	}

	close(deliveryChan)
	return nil
}

func sliceToString(slice []string) string {
	return strings.Join(slice, ",")
}

func (m sql) GetPasswordCredentials(ids []string) ([]types.PasswordCredential, error) {
	var credentials []types.PasswordCredential
	err := m.db.Select(&credentials, "SELECT * FROM password_credentials WHERE id IN ($1)", sliceToString(ids))
	if err != nil {
		return nil, err
	}
	for _, cred := range credentials {
		err := m.produceMessage("creds_read", cred.ID, []byte(fmt.Sprintf("Read PasswordCredential with ID: %s", cred.ID)))
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}
	return credentials, nil
}

func (m sql) GetCardCredentials(ids []string) ([]types.CardCredential, error) {
	var credentials []types.CardCredential
	err := m.db.Select(&credentials, "SELECT * FROM card_credentials WHERE id IN ($1)", sliceToString(ids))
	if err != nil {
		return nil, err
	}
	for _, cred := range credentials {
		err := m.produceMessage("creds_read", cred.ID, []byte(fmt.Sprintf("Read CardCredential with ID: %s", cred.ID)))
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}
	return credentials, nil
}

func (m sql) GetSSHKeyCredentials(ids []string) ([]types.SSHKeyCredential, error) {
	var credentials []types.SSHKeyCredential
	err := m.db.Select(&credentials, "SELECT * FROM ssh_keys WHERE id IN ($1)", sliceToString(ids))
	if err != nil {
		return nil, err
	}
	for _, cred := range credentials {
		err := m.produceMessage("creds_read", cred.ID, []byte(fmt.Sprintf("Read SSHKeyCredential with ID: %s", cred.ID)))
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}
	return credentials, nil
}

func (m sql) CreatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error) {
	var createdCredential types.PasswordCredential
	err := m.db.Get(&createdCredential, "INSERT INTO password_credentials (title, note, user_identifier, password, domain_name) VALUES ($1, $2, $3, $4, $5) RETURNING *", credential.Title, credential.Note, credential.UserIdentifier, credential.Password, credential.DomainName)
	if err != nil {
		return createdCredential, err
	}
	err = m.produceMessage("creds_create", credential.ID, []byte(fmt.Sprintf("Created PasswordCredential with ID: %s", credential.ID)))
	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	return createdCredential, nil
}

func (m sql) CreateCardCredential(credential types.CardCredential) (types.CardCredential, error) {
	var createdCredential types.CardCredential
	err := m.db.Get(&createdCredential ,"INSERT INTO card_credentials (title, note, owner_name, cvc, expiration_date, card_number) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *", credential.Title, credential.Note, credential.OwnerName, credential.CVC, credential.ExpirationDate, credential.CardNumber)
	if err != nil {
		return createdCredential, err
	}
	err = m.produceMessage("creds_create", credential.ID, []byte(fmt.Sprintf("Created CardCredential with ID: %s", credential.ID)))
	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	return createdCredential, nil
}

func (m sql) CreateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error) {
	var createdCredential types.SSHKeyCredential
	err := m.db.Get(&createdCredential, "INSERT INTO ssh_keys (title, note, private_key, public_key, hostname, user_identifier) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *", credential.Title, credential.Note, credential.PrivateKey, credential.PublicKey, credential.Hostname, credential.UserIdentifier)
	if err != nil {
		return createdCredential, err
	}
	err = m.produceMessage("creds_create", credential.ID, []byte(fmt.Sprintf("Created SSHKeyCredential with ID: %s", credential.ID)))
	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	return createdCredential, nil
}

func (m sql) UpdatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error) {
	_, err := m.db.NamedExec("UPDATE password_credentials SET title = :title, note = :note, password = :password, domain_name = :domain_name WHERE id = :id", credential)
	if err != nil {
		return credential, err
	}
	err = m.produceMessage("creds_update", credential.ID, []byte(fmt.Sprintf("Updated PasswordCredential with ID: %s", credential.ID)))
	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	return credential, nil
}

func (m sql) UpdateCardCredential(credential types.CardCredential) (types.CardCredential, error) {
	var updatedCredential types.CardCredential
	_, err := m.db.NamedExec("UPDATE card_credentials SET owner_name = :owner_name, cvc = :cvc, expiration_date = :expiration_date, card_number = :card_number, title = :title, note = :note WHERE id = :id", credential)
	if err != nil {
		return updatedCredential, err
	}
	err = m.produceMessage("creds_update", credential.ID, []byte(fmt.Sprintf("Updated CardCredential with ID: %s", credential.ID)))
	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	return credential, nil
}

func (m sql) UpdateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error) {
	var updatedCredential types.SSHKeyCredential
	_, err := m.db.NamedExec("UPDATE ssh_keys SET private_key = :private_key, public_key = :public_key, hostname = :hostname, user_identifier = :user_identifier, title= :title, note = :note WHERE id = :id", credential)
	if err != nil {
		return updatedCredential, err
	}
	err = m.produceMessage("creds_update", credential.ID, []byte(fmt.Sprintf("Updated SSHKeyCredential with ID: %s", credential.ID)))
	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	return credential, nil
}

func (m sql) DeletePasswordCredentials(ids []string) error {
	_, err := m.db.Exec("DELETE FROM password_credentials WHERE id IN ($1)", sliceToString(ids))
	if err != nil {
		return err
	}
	for _, id := range ids {
		err = m.produceMessage("creds_delete", id, []byte(fmt.Sprintf("Deleted PasswordCredential with ID: %s", id)))
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}
	return nil
}

func (m sql) DeleteCardCredentials(ids []string) error {
	_, err := m.db.Exec("DELETE FROM card_credentials WHERE id IN ($1)", sliceToString(ids))
	if err != nil {
		return err
	}
	for _, id := range ids {
		err = m.produceMessage("creds_delete", id, []byte(fmt.Sprintf("Deleted CardCredential with ID: %s", id)))
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}
	return nil
}

func (m sql) DeleteSSHKeyCredentials(ids []string) error {
	_, err := m.db.Exec("DELETE FROM ssh_keys WHERE id IN ($1)", sliceToString(ids))
	if err != nil {
		return err
	}
	for _, id := range ids {
		err = m.produceMessage("creds_delete", id, []byte(fmt.Sprintf("Deleted SSHKeyCredential with ID: %s", id)))
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}
	return nil
}