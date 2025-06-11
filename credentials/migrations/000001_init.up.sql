CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS credentials (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  title VARCHAR(1000),
  note TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP,
  last_read_at TIMESTAMP,
  custom_fields JSONB
);

CREATE TABLE IF NOT EXISTS card_credentials (
  owner_name VARCHAR(25) NOT NULL,
  CVC INTEGER NOT NULL,
  expiration_date VARCHAR(25) NOT NULL,
  card_number NUMERIC(16,0) NOT NULL
) INHERITS (credentials);

CREATE TABLE IF NOT EXISTS password_credentials (
  user_identifier VARCHAR(1000) NOT NULL,
  password VARCHAR(65534) NOT NULL,
  domain_name VARCHAR(50000) NOT NULL
) INHERITS (credentials);

CREATE TABLE IF NOT EXISTS ssh_keys (
  private_key TEXT NOT NULL,
  public_key TEXT NOT NULL,
  hostname VARCHAR(500),
  user_identifier VARCHAR(500)
) INHERITS (credentials);
