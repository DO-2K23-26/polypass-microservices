package organizations

type Config struct {
	Endpoint string `mapstructure:"endpoint" env:"ORGANIZATIONS_ENDPOINT"`
}
