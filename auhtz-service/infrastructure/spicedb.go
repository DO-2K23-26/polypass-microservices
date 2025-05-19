package infrastructure

import (
	"context"
	"log"
	"os"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/common/types"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
)

type SpiceDBAdapter struct {
	authzedClient *authzed.Client
}

func NewSpiceDBAdapter(authzedClient *authzed.Client) *SpiceDBAdapter {
	return &SpiceDBAdapter{
		authzedClient: authzedClient,
	}
}

func (s *SpiceDBAdapter) Close() error {
	return s.authzedClient.Close()
}

func (s *SpiceDBAdapter) HealthCheck() error {
	// Example health check: Attempt to read relationships
	ctx := context.Background()
	_, err := s.authzedClient.ReadSchema(ctx, &v1.ReadSchemaRequest{})
	if err != nil {
		log.Println("health check failed: %w", err)
		return err
	}
	return nil
}

func (s *SpiceDBAdapter) Migrate() error {
	schema, err := os.ReadFile("schema.zed")
	if err != nil {
		log.Printf("failed to read schema file: %v", err)
		return err
	}
	schemaString := string(schema)

	// Check if schema is up to date
	isUpToDate, err := s.authzedClient.SchemaServiceClient.DiffSchema(context.Background(), &v1.DiffSchemaRequest{
		ComparisonSchema: schemaString,
		Consistency: &v1.Consistency{
			Requirement: &v1.Consistency_FullyConsistent{FullyConsistent: true},
		},
	})
	if err != nil {
		log.Printf("failed to diff schema: %v", err)
		return err
	}
	if len(isUpToDate.Diffs) == 0 {
		log.Println("Authzed schema is up to date")
		return nil
	}
	res, err := s.authzedClient.SchemaServiceClient.WriteSchema(context.Background(), &v1.WriteSchemaRequest{
		Schema: schemaString,
	})
	if err != nil {
		log.Println("Could not write the schema to spicedb:", err)
		return err
	}
	log.Println("Wrote the schema to spicedb", res)

	return nil
}

func (s *SpiceDBAdapter) CreateRelationship(
	ctx context.Context,
	subjectType types.ObjectType,
	subjectId string,
	relation types.Relation,
	resourceType types.ObjectType,
	resourceId string,
) error {
	_, err := s.authzedClient.WriteRelationships(ctx, &v1.WriteRelationshipsRequest{
		Updates: []*v1.RelationshipUpdate{
			{
				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
				Relationship: &v1.Relationship{
					Subject: &v1.SubjectReference{
						Object: &v1.ObjectReference{ObjectType: string(subjectType), ObjectId: subjectId},
					},
					Resource: &v1.ObjectReference{ObjectType: string(resourceType), ObjectId: resourceId},
					Relation: string(relation),
				},
			},
		},
	})
	if err != nil {
		log.Printf("failed to create relationship: %v", err)
		return err
	}
	return nil
}

func (s *SpiceDBAdapter) UpdateRelationship(
	ctx context.Context,
	subjectType types.ObjectType,
	subjectId string,
	relation types.Relation,
	resourceType types.ObjectType,
	resourceId string,
) error {
	// Perform the delete and create operations in a single transaction
	_, err := s.authzedClient.WriteRelationships(ctx, &v1.WriteRelationshipsRequest{
		Updates: []*v1.RelationshipUpdate{
			// Delete all existing relationships for the given resource and relation
			{
				Operation: v1.RelationshipUpdate_OPERATION_DELETE,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{
						ObjectType: string(resourceType),
						ObjectId:   resourceId,
					},
					Relation: string(relation),
				},
			},
			// Create the new relationship
			{
				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
				Relationship: &v1.Relationship{
					Subject: &v1.SubjectReference{
						Object: &v1.ObjectReference{ObjectType: string(subjectType), ObjectId: subjectId},
					},
					Resource: &v1.ObjectReference{ObjectType: string(resourceType), ObjectId: resourceId},
					Relation: string(relation),
				},
			},
		},
	})
	if err != nil {
		log.Printf("failed to update relationship: %v", err)
		return err
	}

	log.Printf("successfully updated relationship between subject %s:%s and resource %s:%s with relation %s", subjectType, subjectId, resourceType, resourceId, relation)
	return nil
}

func (s *SpiceDBAdapter) DeleteRelationship(
	ctx context.Context,
	subjectType types.ObjectType,
	subjectId string,
	resourceType types.ObjectType,
	resourceId string,
) error {
	_, err := s.authzedClient.DeleteRelationships(ctx, &v1.DeleteRelationshipsRequest{
		RelationshipFilter: &v1.RelationshipFilter{
			ResourceType:       string(resourceType),
			OptionalResourceId: resourceId,
			OptionalSubjectFilter: &v1.SubjectFilter{
				SubjectType:       string(subjectType),
				OptionalSubjectId: subjectId,
			},
		},
	})
	if err != nil {
		log.Printf("failed to delete relationship: %v", err)
		return err
	}
	return nil
}

func (s *SpiceDBAdapter) Delete(ctx context.Context, resourceType types.ObjectType, id string) error {
	_, err := s.authzedClient.DeleteRelationships(ctx, &v1.DeleteRelationshipsRequest{

		RelationshipFilter: &v1.RelationshipFilter{
			OptionalResourceId: id,
			ResourceType:       string(resourceType),
		},
	})
	if err != nil {
		log.Printf("failed to delete all occurrences of id %s: %v", id, err)
		return err
	}

	log.Printf("successfully deleted all occurrences of id %s", id)
	return nil
}

func (s *SpiceDBAdapter) ChangeParent(ctx context.Context, folderId string, newParentId string) error {
	// Perform the delete and create operations in a single transaction
	_, err := s.authzedClient.WriteRelationships(ctx, &v1.WriteRelationshipsRequest{
		Updates: []*v1.RelationshipUpdate{
			// Delete all existing parent relationships for the folder
			{
				Operation: v1.RelationshipUpdate_OPERATION_DELETE,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{
						ObjectType: "folder",
						ObjectId:   folderId,
					},
					Relation: "parent",
				},
			},
			// Create the new parent relationship
			{
				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{
						ObjectType: "folder",
						ObjectId:   folderId,
					},
					Relation: "parent",
					Subject: &v1.SubjectReference{
						Object: &v1.ObjectReference{
							ObjectType: "folder",
							ObjectId:   newParentId,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("failed to change parent relationship: %v", err)
		return err
	}

	log.Printf("successfully changed parent of folder %s to %s", folderId, newParentId)
	return nil
}
