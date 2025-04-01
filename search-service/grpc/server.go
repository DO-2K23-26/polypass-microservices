package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/DO-2K23-26/polypass-microservices/search-service/gen/search/api"
	"github.com/DO-2K23-26/polypass-microservices/search-service/services/credential"
	"github.com/DO-2K23-26/polypass-microservices/search-service/services/folder"
	tag "github.com/DO-2K23-26/polypass-microservices/search-service/services/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the gRPC server
type Server struct {
	server   *grpc.Server
	listener net.Listener
}

// NewServer creates a new gRPC server with the necessary services
func NewServer(
	credentialService *credential.CredentialService,
	folderService *folder.FolderService,
	tagService *tag.TagService,
	port int,
) (*Server, error) {
	// Create a listener on specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Create the search service implementation
	searchService := NewSearchServiceServer(
		credentialService,
		folderService,
		tagService,
	)

	api.RegisterSearchServiceServer(grpcServer, searchService)

	// Register reflection service (useful for debugging with tools like grpcurl)
	reflection.Register(grpcServer)

	return &Server{
		server:   grpcServer,
		listener: lis,
	}, nil
}

// Start starts the gRPC server
func (s *Server) Start() error {
	log.Printf("gRPC server starting on %s", s.listener.Addr())
	return s.server.Serve(s.listener)
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	s.server.GracefulStop()
}
