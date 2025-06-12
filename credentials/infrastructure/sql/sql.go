package sql

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DO-2K23-26/polypass-microservices/credentials/types"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/linkedin/goavro"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

import (
	avro "github.com/DO-2K23-26/polypass-microservices/credentials/interfaces/credentials"
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
	schemaDir  string
}

func NewSql(config Config, producer *kafka.Producer, consumer *kafka.Consumer) (Sql, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Dbname))
	if err != nil {
		return nil, err
	}
	return sql{
		db:         db,
		producer:   producer,
		consumer:   consumer,
		migrations: config.Migrations,
		username:   config.Username,
		password:   config.Password,
		host:       config.Host,
		port:       config.Port,
		dbname:     config.Dbname,
		schemaDir:  config.SchemaDir,
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
	"PasswordCredential":       "password_credential.avsc",
    "CardCredential":           "card_credential.avsc",
    "SSHKeyCredential":         "ssh_credential.avsc",
    "CardAttribute":            "card_attributes.avsc",
    "CreateCredentialsOpts":    "create_credentials_opts.avsc",
    "Credential":               "credential.avsc",
    "PasswordAttribute":        "password_attributes.avsc",
    "SSHKeyAttribute":          "ssh_attribute.avsc",
    "UserIdentifierAttributes": "user_identifier_attributes.avsc",
}

func (m *sql) loadSchema(filename string) (string, error) {
	data, err := avro.FS.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("lecture schema %s: %w", filename, err)
	}
	return string(data), nil
}

func getSchemaPath(typeName string) (string, error) {
	path, ok := credentialsSchema[typeName]
	if !ok || path == "" {
		return "", fmt.Errorf("no schema found for type %s", typeName)
	}
	return path, nil
}

// newProduceMessage prend n'importe quel credential et le sérialise
func (m *sql) produceMessage(topic string, cred interface{}) error {
	var (
		typeName string
		record   map[string]interface{}
	)
	switch c := cred.(type) {
	case types.PasswordCredential:
		typeName = "PasswordCredential"
		record = map[string]interface{}{
			"Credential": map[string]interface{}{
				"id":            c.Credential.ID,
				"title":         c.Credential.Title,
				"note":          c.Credential.Note,
				"created_at":    unixOrZero(c.Credential.CreatedAt),
				"updated_at":    unixOrZero(c.Credential.UpdatedAt),
				"expires_at":    unixOrZero(c.Credential.ExpiresAt),
				"last_read_at":  unixOrZero(c.Credential.LastReadAt),
				"custom_fields": toInterfaceMap(c.Credential.CustomFields),
			},
			"PasswordAttributes": map[string]interface{}{
				"password":    c.PasswordAttributes.Password,
				"domain_name": c.PasswordAttributes.DomainName,
			},
			"UserIdentifierAttribute": map[string]interface{}{
				"user_identifier": c.UserIdentifierAttribute.UserIdentifier,
			},
		}

	case types.CardCredential:
		typeName = "CardCredential"
		record = map[string]interface{}{
			"Credential": map[string]interface{}{
				"id":            c.Credential.ID,
				"title":         c.Credential.Title,
				"note":          c.Credential.Note,
				"created_at":    unixOrZero(c.Credential.CreatedAt),
				"updated_at":    unixOrZero(c.Credential.UpdatedAt),
				"expires_at":    unixOrZero(c.Credential.ExpiresAt),
				"last_read_at":  unixOrZero(c.Credential.LastReadAt),
				"custom_fields": toInterfaceMap(c.Credential.CustomFields),
			},
			"CardAttributes": map[string]interface{}{
				"owner_name":      c.OwnerName,
				"cvc":             c.CVC,
				"expiration_date": c.ExpirationDate,
				"card_number":     c.CardNumber,
			},
			"UserIdentifierAttribute": map[string]interface{}{
				"user_identifier": c.UserIdentifier,
			},
		}

	case types.SSHKeyCredential:
		typeName = "SSHKeyCredential"
		record = map[string]interface{}{
			"Credential": map[string]interface{}{
                "id":            c.Credential.ID,
                "title":         c.Credential.Title,
                "note":          c.Credential.Note,
                "created_at":    unixOrZero(c.Credential.CreatedAt),
                "updated_at":    unixOrZero(c.Credential.UpdatedAt),
                "expires_at":    unixOrZero(c.Credential.ExpiresAt),
                "last_read_at":  unixOrZero(c.Credential.LastReadAt),
                "custom_fields": toInterfaceMap(c.Credential.CustomFields),
            },
			"SSHKeyAttributes": map[string]interface{}{
				"private_key":     c.PrivateKey,
				"public_key":      c.PublicKey,
				"hostname":        c.Hostname,
				"user_identifier": c.UserIdentifier,
			},
			"UserIdentifierAttribute": map[string]interface{}{
				"user_identifier": c.UserIdentifier,
			},
		}
	default:
		return fmt.Errorf("unsupported credential type %T", cred)
	}

	schemaPath, err := getSchemaPath(typeName)
	if err != nil {
		return err
	}
	schemaDef, err := m.loadSchema(schemaPath)
	if err != nil {
		return err
	}
	codec, err := goavro.NewCodec(schemaDef)
	if err != nil {
		return err
	}

	avroBin, err := codec.BinaryFromNative(nil, record)
	if err != nil {
		log.Printf("Failed to serialize data: %v", err)
		return err
	}

	dc := make(chan kafka.Event, 1)
	err = m.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(typeName),
		Value:          avroBin,
	}, dc)
	if err != nil {
		return err
	}
	ev := <-dc
	if km := ev.(*kafka.Message); km.TopicPartition.Error != nil {
		return km.TopicPartition.Error
	}
	close(dc)
	log.Printf("[INFO] Message produit sur Kafka — topic: %s | key: %s | type: %s | record=%+v", topic, typeName, typeName, record)
	return nil
}

func unixOrZero(t *time.Time) int64 {
	if t != nil {
		return t.Unix()
	}
	return 0
}

func toInterfaceMap(m *map[string]any) map[string]interface{} {
	out := make(map[string]interface{})
	if m == nil {
		return out
	}
	for k, v := range *m {
		out[k] = v
	}
	return out
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
		err := m.produceMessage("creds_read", cred)
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
		err := m.produceMessage("creds_read", cred)
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
		err := m.produceMessage("creds_read", cred)
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
	err = m.produceMessage("creds_create", createdCredential)

	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	return createdCredential, nil
}

func (m sql) CreateCardCredential(credential types.CardCredential) (types.CardCredential, error) {
	var createdCredential types.CardCredential
	err := m.db.Get(&createdCredential, "INSERT INTO card_credentials (title, note, owner_name, cvc, expiration_date, card_number) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *", credential.Title, credential.Note, credential.OwnerName, credential.CVC, credential.ExpirationDate, credential.CardNumber)
	if err != nil {
		return createdCredential, err
	}
	err = m.produceMessage("creds_create", createdCredential)
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
	err = m.produceMessage("creds_create", createdCredential)
	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	return createdCredential, nil
}

func (m sql) UpdatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error) {
	now := time.Now()
    credential.UpdatedAt = &now

    _, err := m.db.NamedExec(`
        UPDATE password_credentials
        SET title       = :title,
            note        = :note,
            password    = :password,
            domain_name = :domain_name,
            updated_at  = :updated_at
        WHERE id = :id
    `, credential)
    if err != nil {
        return credential, err
    }

    if err := m.produceMessage("creds_update", credential); err != nil {
        log.Printf("Failed to produce message: %v", err)
    }
    return credential, nil
}

func (m sql) UpdateCardCredential(credential types.CardCredential) (types.CardCredential, error) {
	now := time.Now()
	credential.UpdatedAt = &now
	_, err := m.db.NamedExec(`
      UPDATE card_credentials
      SET owner_name      = :owner_name,
          cvc             = :cvc,
          expiration_date = :expiration_date,
          card_number     = :card_number,
          title           = :title,
          note            = :note,
          updated_at      = :updated_at
      WHERE id = :id
    `, credential)
	if err != nil {
		return credential, err
	}
	if err := m.produceMessage("creds_update", credential); err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
	return credential, nil
}

func (m sql) UpdateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error) {
	now := time.Now()
	credential.UpdatedAt = &now

    _, err := m.db.NamedExec(`
        UPDATE ssh_keys
        SET private_key    = :private_key,
            public_key     = :public_key,
            hostname       = :hostname,
            user_identifier= :user_identifier,
            title          = :title,
            note           = :note,
            updated_at     = :updated_at
        WHERE id = :id
    `, credential)
    if err != nil {
        return credential, err
    }

    if err := m.produceMessage("creds_update", credential); err != nil {
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
		err = m.produceMessage("creds_delete", id)
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
		err = m.produceMessage("creds_delete", id)
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
		err = m.produceMessage("creds_delete", id)
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}
	return nil
}
