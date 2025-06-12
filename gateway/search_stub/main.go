// filepath: gateway/search_stub/main.go
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/DO-2K23-26/polypass-microservices/gateway/proto/search"
)

// mockServer implements the gRPC SearchService with dummy data
type mockServer struct {
	pb.UnimplementedSearchServiceServer
}

func (s *mockServer) SearchFolders(ctx context.Context, req *pb.SearchFoldersRequest) (*pb.SearchFoldersResponse, error) {
	return &pb.SearchFoldersResponse{
		Total: 1,
		Folders: []*pb.Folder{{
			Id:       "f1",
			Name:     "MockFolder",
			ParentId: "",
		}},
	}, nil
}

func (s *mockServer) SearchTags(ctx context.Context, req *pb.SearchTagsRequest) (*pb.SearchTagsResponse, error) {
	return &pb.SearchTagsResponse{
		Total: 1,
		Tags: []*pb.Tag{{
			Id:       "t1",
			Name:     "MockTag",
			FolderId: "f1",
		}},
	}, nil
}

func (s *mockServer) SearchCredentials(ctx context.Context, req *pb.SearchCredentialsRequest) (*pb.SearchCredentialsResponse, error) {
	return &pb.SearchCredentialsResponse{
		Total: 1,
		Credentials: []*pb.Credential{{
			Id:       "c1",
			Title:    "MockCredential",
			FolderId: "f1",
			Folder: &pb.Folder{Id: "f1", Name: "MockFolder", ParentId: ""},
			Tags: []*pb.Tag{{Id: "t1", Name: "MockTag", FolderId: "f1"}},
		}},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSearchServiceServer(grpcServer, &mockServer{})
	log.Println("Mock SearchService gRPC running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
