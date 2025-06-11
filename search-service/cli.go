package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    "github.com/DO-2K23-26/polypass-microservices/search-service/gen/search/api"
)

func main() {
    // Initialize Viper for configuration
    viper.SetDefault("grpc_server", "localhost:8081")
    viper.AutomaticEnv()

    // Create the root command
    var rootCmd = &cobra.Command{
        Use:   "search-cli",
        Short: "CLI for testing the SearchService gRPC API",
    }

    // Add subcommands
    rootCmd.AddCommand(searchFoldersCmd())
    rootCmd.AddCommand(searchTagsCmd())
    rootCmd.AddCommand(searchCredentialsCmd())

    // Execute the CLI
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func getClient() (api.SearchServiceClient, *grpc.ClientConn) {
    // Connect to the gRPC server
    conn, err := grpc.Dial(viper.GetString("grpc_server"), grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }

    // Create the client
    client := api.NewSearchServiceClient(conn)
    return client, conn
}

func searchFoldersCmd() *cobra.Command {
    var searchQuery string
    var limit, page int32
    var userID string

    cmd := &cobra.Command{
        Use:   "search-folders",
        Short: "Search for folders",
        Run: func(cmd *cobra.Command, args []string) {
            client, conn := getClient()
            defer conn.Close()

            // Create the request
            req := &api.SearchFoldersRequest{
                SearchQuery: searchQuery,
                Limit:       limit,
                Page:        page,
                UserId:      userID,
            }

            // Call the gRPC method
            ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
            defer cancel()

            resp, err := client.SearchFolders(ctx, req)
            if err != nil {
                log.Fatalf("Error calling SearchFolders: %v", err)
            }

            fmt.Printf("Folders: %+v\n", resp.Folders)
            fmt.Printf("Total: %d\n", resp.Total)
        },
    }

    cmd.Flags().StringVar(&searchQuery, "query", "", "Search query")
    cmd.Flags().Int32Var(&limit, "limit", 10, "Limit")
    cmd.Flags().Int32Var(&page, "page", 1, "Page")
    cmd.Flags().StringVar(&userID, "user-id", "", "User ID")
    return cmd
}

func searchTagsCmd() *cobra.Command {
    var searchQuery, folderID, userID string
    var limit, page int32

    cmd := &cobra.Command{
        Use:   "search-tags",
        Short: "Search for tags",
        Run: func(cmd *cobra.Command, args []string) {
            client, conn := getClient()
            defer conn.Close()

            // Create the request
            req := &api.SearchTagsRequest{
                SearchQuery: searchQuery,
                FolderId:    folderID,
                Limit:       limit,
                Page:        page,
                UserId:      userID,
            }

            // Call the gRPC method
            ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
            defer cancel()

            resp, err := client.SearchTags(ctx, req)
            if err != nil {
                log.Fatalf("Error calling SearchTags: %v", err)
            }

            fmt.Printf("Tags: %+v\n", resp.Tags)
            fmt.Printf("Total: %d\n", resp.Total)
        },
    }

    cmd.Flags().StringVar(&searchQuery, "query", "", "Search query")
    cmd.Flags().StringVar(&folderID, "folder-id", "", "Folder ID")
    cmd.Flags().Int32Var(&limit, "limit", 10, "Limit")
    cmd.Flags().Int32Var(&page, "page", 1, "Page")
    cmd.Flags().StringVar(&userID, "user-id", "", "User ID")
    return cmd
}

func searchCredentialsCmd() *cobra.Command {
    var searchQuery, folderID, userID string
    var tagIDs []string
    var limit, page int32

    cmd := &cobra.Command{
        Use:   "search-credentials",
        Short: "Search for credentials",
        Run: func(cmd *cobra.Command, args []string) {
            client, conn := getClient()
            defer conn.Close()

            // Create the request
            req := &api.SearchCredentialsRequest{
                SearchQuery: searchQuery,
                FolderId:    folderID,
                TagIds:      tagIDs,
                Limit:       limit,
                Page:        page,
                UserId:      userID,
            }

            // Call the gRPC method
            ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
            defer cancel()

            resp, err := client.SearchCredentials(ctx, req)
            if err != nil {
                log.Fatalf("Error calling SearchCredentials: %v", err)
            }

            fmt.Printf("Credentials: %+v\n", resp.Credentials)
            fmt.Printf("Total: %d\n", resp.Total)
        },
    }

    cmd.Flags().StringVar(&searchQuery, "query", "", "Search query")
    cmd.Flags().StringVar(&folderID, "folder-id", "", "Folder ID")
    cmd.Flags().StringSliceVar(&tagIDs, "tag-ids", []string{}, "Tag IDs")
    cmd.Flags().Int32Var(&limit, "limit", 10, "Limit")
    cmd.Flags().Int32Var(&page, "page", 1, "Page")
    cmd.Flags().StringVar(&userID, "user-id", "", "User ID")
    return cmd
}