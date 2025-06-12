package graphql

type Config struct {
	ListenAddr string `mapstructure:"listen_addr" env:"HTTP_LISTEN_ADDR"`
}
