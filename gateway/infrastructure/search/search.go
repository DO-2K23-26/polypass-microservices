// filepath: gateway/infrastructure/search/search.go
package search

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"github.com/DO-2K23-26/polypass-microservices/gateway/proto/search"
	"github.com/DO-2K23-26/polypass-microservices/gateway/infrastructure"
)

// SearchAPI defines the methods to interact with the gRPC SearchService
type SearchAPI interface {
	infrastructure.Repository
	SearchServiceClient() search.SearchServiceClient
}

// client implements SearchAPI and repository
type client struct {
	cfg    Config
	conn   *grpc.ClientConn
	client search.SearchServiceClient
}

// New creates a new SearchAPI client
func New(cfg Config) SearchAPI {
	return &client{cfg: cfg}
}

// Setup establishes the gRPC connection
func (c *client) Setup() error {
	conn, err := grpc.Dial(c.cfg.Endpoint, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("failed to dial search service: %w", err)
	}
	c.conn = conn
	c.client = search.NewSearchServiceClient(conn)
	return nil
}

// Shutdown closes the gRPC connection
func (c *client) Shutdown() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// SearchServiceClient returns the underlying gRPC client
func (c *client) SearchServiceClient() search.SearchServiceClient {
	return c.client
}
