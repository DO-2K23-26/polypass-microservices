package repositories

import (
	"context"
	"os"
	"sharing/models"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StoredSecret struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Content      map[string]string `bson:"content"`
	Expiration   primitive.DateTime `bson:"expiration"`
	IsEncrypted  bool              `bson:"is_encrypted"`
	IsOneTimeUse bool              `bson:"is_one_time_use"`
}

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
		collection := db.Collection("secrets")
		// Create an index on the "expiration" field
		indexModel := mongo.IndexModel{
			Keys: bson.D{
				{Key: "expiration", Value: 1},
			},
			Options: options.Index().SetExpireAfterSeconds(1),
		}

		_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)

		if err != nil {
			panic(err)
		}
    return &secretRepository{db: db}
}

type secretRepository struct {
    db *mongo.Database
}

func (r *secretRepository) CreateSecret(secret *models.Secret) (models.Secret, error) {
	collection := r.db.Collection("secrets")

	storedSecret := StoredSecret{
		Content:      secret.Content,
		Expiration:   primitive.NewDateTimeFromTime(time.Unix(secret.Expiration, 0)),
		IsEncrypted:  secret.IsEncrypted,
		IsOneTimeUse: secret.IsOneTimeUse,
	}

	// Insert the secret into the collection
	result, err := collection.InsertOne(context.TODO(), storedSecret)
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
	var storedSecret StoredSecret
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
			return nil, err
	}

	err = collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&storedSecret)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Return nil if the document has expired and is no longer available
	} else if err != nil {
		return nil, err
	}

	secret = models.Secret{
		Id:           storedSecret.Id.Hex(),
		Content:      storedSecret.Content,
		Expiration:   storedSecret.Expiration.Time().Unix(),
		IsEncrypted:  storedSecret.IsEncrypted,
		IsOneTimeUse: storedSecret.IsOneTimeUse,
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