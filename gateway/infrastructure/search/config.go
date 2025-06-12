// filepath: gateway/infrastructure/search/config.go
package search

type Config struct {
	// gRPC endpoint for SearchService, e.g., "localhost:50051"
	Endpoint string `mapstructure:"endpoint" env:"SEARCH_ENDPOINT"`
}
