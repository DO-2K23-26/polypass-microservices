package sql

import (
	"fmt"

	"github.com/DO-2K23-26/polypass-microservices/credentials/types"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"

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
	migrations string
	username   string
	password   string
	host       string
	port       int
	dbname     string
}

func NewSql(config Config) (Sql, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Dbname))
	if err != nil {
		return nil, err
	}
	return sql{
		db:         db,
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

func (m sql) GetPasswordCredentials(ids []string) ([]types.PasswordCredential, error) {
	var credentials []types.PasswordCredential
	err := m.db.Select(&credentials, "SELECT * FROM password_credentials WHERE id = ANY($1)", ids)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func (m sql) GetCardCredentials(ids []string) ([]types.CardCredential, error) {
	var credentials []types.CardCredential
	err := m.db.Select(&credentials, "SELECT * FROM card_credentials WHERE id = ANY($1)", ids)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func (m sql) GetSSHKeyCredentials(ids []string) ([]types.SSHKeyCredential, error) {
	var credentials []types.SSHKeyCredential
	err := m.db.Select(&credentials, "SELECT * FROM ssh_keys WHERE id = ANY($1)", ids)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func (m sql) CreatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error) {
	var createdCredential types.PasswordCredential
	_, err := m.db.NamedExec("INSERT INTO password_credentials (title, note, user_identifier, password, domain_name) VALUES (:user_identifier, :password, :domain_name, :title, :note)", credential)
	if err != nil {
		return createdCredential, err
	}
	return credential, nil
}

func (m sql) CreateCardCredential(credential types.CardCredential) (types.CardCredential, error) {
	var createdCredential types.CardCredential
	_, err := m.db.NamedExec("INSERT INTO card_credentials (title, note, owner_name, cvc, expiration_date, card_number) VALUES (:title, :note, :owner_name, :cvc, :expiration_date, :card_number)", credential)
	if err != nil {
		return createdCredential, err
	}
	return credential, nil
}

func (m sql) CreateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error) {
	var createdCredential types.SSHKeyCredential
	_, err := m.db.NamedExec("INSERT INTO ssh_keys (title, note, private_key, public_key, hostname, user_identifier) VALUES (:title, :note, :private_key, :public_key, :hostname, :user_identifier)", credential)
	if err != nil {
		return createdCredential, err
	}
	return credential, nil
}

func (m sql) UpdatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error) {
	var updatedCredential types.PasswordCredential
	_, err := m.db.NamedExec("UPDATE password_credentials SET title = :title, note = :note, password = :password, domain_name = :domain_name WHERE id = :id", credential)
	if err != nil {
		return updatedCredential, err
	}
	return credential, nil
}

func (m sql) UpdateCardCredential(credential types.CardCredential) (types.CardCredential, error) {
	var updatedCredential types.CardCredential
	_, err := m.db.NamedExec("UPDATE card_credentials SET owner_name = :owner_name, cvc = :cvc, expiration_date = :expiration_date, card_number = :card_number, title = :title, note = :note WHERE id = :id", credential)
	if err != nil {
		return updatedCredential, err
	}
	return credential, nil
}

func (m sql) UpdateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error) {
	var updatedCredential types.SSHKeyCredential
	_, err := m.db.NamedExec("UPDATE ssh_keys SET private_key = :private_key, public_key = :public_key, hostname = :hostname, user_identifier = :user_identifier, title= :title, note = :note WHERE id = :id", credential)
	if err != nil {
		return updatedCredential, err
	}
	return credential, nil
}

func (m sql) DeletePasswordCredentials(ids []string) error {
	_, err := m.db.Exec("DELETE FROM password_credentials WHERE id = ANY($1)", ids)
	if err != nil {
		return err
	}
	return nil
}

func (m sql) DeleteCardCredentials(ids []string) error {
	_, err := m.db.Exec("DELETE FROM card_credentials WHERE id = ANY($1)", ids)
	if err != nil {
		return err
	}
	return nil
}

func (m sql) DeleteSSHKeyCredentials(ids []string) error {
	_, err := m.db.Exec("DELETE FROM ssh_keys WHERE id = ANY($1)", ids)
	if err != nil {
		return err
	}
	return nil
}
