package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"
    "encoding/json"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    "github.com/hamba/avro"
    "github.com/confluentinc/confluent-kafka-go/kafka"

    "github.com/DO-2K23-26/polypass-microservices/search-service/gen/search/api"
)

func main() {
    // Initialize Viper for configuration
    viper.SetDefault("grpc_server", "localhost:8081")
    viper.AutomaticEnv()

    // Create the root command
    var rootCmd = &cobra.Command{
        Use:   "search-cli",
        Short: "CLI for testing the SearchService gRPC API and Kafka",
    }

    // Add subcommands
    rootCmd.AddCommand(kafkaProducerCmd())
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

var tagCreationSchema = `
{
    "type": "record",
    "name": "TagCreation",
    "fields": [
        {"name": "id", "type": "string"},
        {"name": "name", "type": "string"},
        {"name": "folder_id", "type": "string"}
    ]
}`

func kafkaProducerCmd() *cobra.Command {
    var topic, schemaType, message string

    cmd := &cobra.Command{
        Use:   "produce-kafka",
        Short: "Produce Kafka messages in Avro format",
        Run: func(cmd *cobra.Command, args []string) {
            producer, err := kafka.NewProducer(&kafka.ConfigMap{
                "bootstrap.servers": viper.GetString("kafka_server"),
            })
            if err != nil {
                log.Fatalf("Failed to create Kafka producer: %v", err)
            }
            defer func() {
                // Ensure all messages are delivered before exiting
                producer.Flush(5000) // Wait up to 5 seconds for delivery
                producer.Close()
            }()

            // Select schema based on schemaType
            var schema avro.Schema
            switch schemaType {
            case "tag_creation":
                schema, err = avro.Parse(tagCreationSchema)
                if err != nil {
                    log.Fatalf("Failed to parse Avro schema: %v", err)
                }
            default:
                log.Fatalf("Unsupported schema type: %s", schemaType)
            }

            // Parse the message into a map
            var messageData map[string]interface{}
            err = json.Unmarshal([]byte(message), &messageData)
            if err != nil {
                log.Fatalf("Failed to parse message JSON: %v", err)
            }

            // Serialize the message using Avro
            encodedMessage, err := avro.Marshal(schema, messageData)
            if err != nil {
                log.Fatalf("Failed to encode message: %v", err)
            }

            // Produce the message to Kafka
            err = producer.Produce(&kafka.Message{
                TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
                Value:          encodedMessage,
            }, nil)
            if err != nil {
                log.Fatalf("Failed to produce message: %v", err)
            }

            fmt.Printf("Message produced to topic %s\n", topic)
        },
    }

    cmd.Flags().StringVar(&topic, "topic", "", "Kafka topic to produce to")
    cmd.Flags().StringVar(&schemaType, "schema", "", "Schema type (e.g., tag_creation)")
    cmd.Flags().StringVar(&message, "message", "", "Message JSON to produce")
    cmd.MarkFlagRequired("topic")
    cmd.MarkFlagRequired("schema")
    cmd.MarkFlagRequired("message")

    return cmd
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
