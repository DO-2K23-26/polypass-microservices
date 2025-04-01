package repositories

import (
	"context"
	"os"
	"sharing/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SecretRepository interface {
    CreateSecret(secret *models.Secret) (models.Secret, error)
    GetSecret(id string) (*models.Secret, error)
}

func NewSecretRepository() SecretRepository {
    // Load environment variables
    godotenv.Load()

    mongoURI := os.Getenv("MONGO_URI")
    clientOptions := options.Client().ApplyURI(mongoURI)

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        panic(err)
    }

    // Verify connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        panic(err)
    }

    // Return the repository
    db := client.Database(os.Getenv("MONGO_DB_NAME"))
    return &secretRepository{db: db}
}

type secretRepository struct {
    db *mongo.Database
}

func (r *secretRepository) CreateSecret(secret *models.Secret) (models.Secret, error) {
	collection := r.db.Collection("secrets")

	// Insert the secret into the collection
	result, err := collection.InsertOne(context.TODO(), secret)
	if err != nil {
		return models.Secret{}, err
	}

	// Set the generated ID back to the secret
    secret.Id = result.InsertedID.(primitive.ObjectID).Hex()

	return *secret, nil
}

func (r *secretRepository) GetSecret(id string) (*models.Secret, error) {
    collection := r.db.Collection("secrets")

    // Find the secret by ID
    var secret models.Secret
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    err = collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&secret)
    if err != nil {
        return nil, err
    }

		if secret.IsOneTimeUse {
			// Delete the secret after it has been accessed
			_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
			if err != nil {
				return nil, err
			}
		}

    return &secret, nil
}