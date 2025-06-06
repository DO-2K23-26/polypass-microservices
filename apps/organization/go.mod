module github.com/DO-2K23-26/polypass-microservices/organization

go 1.24.3

require (
	github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas v0.0.0
	github.com/DO-2K23-26/polypass-microservices/libs/interfaces v0.0.0
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/gorilla/mux v1.8.1
	github.com/spf13/viper v1.20.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/actgardner/gogen-avro/v10 v10.2.1 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/linkedin/goavro/v2 v2.13.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/riferrei/srclient v0.7.2 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.0.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.6.0 // indirect
	gorm.io/driver/postgres v1.6.0 // indirect
	gorm.io/gorm v1.30.0 // indirect
)

replace github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas => ../../libs/avro-schemas

replace github.com/DO-2K23-26/polypass-microservices/libs/interfaces => ../../libs/interfaces
