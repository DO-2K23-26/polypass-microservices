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
	Name		 string            `bson:"name"`
	Content      map[string]string `bson:"content"`
	Expiration   primitive.DateTime `bson:"expiration"`
	IsEncrypted  bool              `bson:"is_encrypted"`
	IsOneTimeUse bool              `bson:"is_one_time_use"`
	User         string				`bson:"user_id"`
	CreatedAt    primitive.DateTime `bson:"created_at"`
}

type SecretRepository interface {
    CreateSecret(secret *models.Secret) (models.Secret, error)
    GetSecret(id string) (*models.Secret, error)
	GetHistory(userId string) ([]models.HistorySecret, error)
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
		Name:         secret.Name,
		Expiration:   primitive.NewDateTimeFromTime(time.Unix(secret.Expiration, 0)),
		IsEncrypted:  secret.IsEncrypted,
		IsOneTimeUse: secret.IsOneTimeUse,
		User: 		  secret.User,
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
		Name:         storedSecret.Name,
		CreatedAt:    storedSecret.CreatedAt.Time().Unix(),
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

func (r *secretRepository) GetHistory(userId string) ([]models.HistorySecret, error) {
	collection := r.db.Collection("secrets")

	// Find all secrets for the given user_id
	filter := bson.M{"user_id": userId}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var history []models.HistorySecret
	for cursor.Next(context.TODO()) {
		var storedSecret StoredSecret
		if err := cursor.Decode(&storedSecret); err != nil {
			return nil, err
		}

		history = append(history, models.HistorySecret{
			Id:           storedSecret.Id.Hex(),
			CreatedAt:    storedSecret.CreatedAt.Time().Unix(),
			Expiration:   storedSecret.Expiration.Time().Unix(),
			ContentSize:  len(storedSecret.Content),
			IsOneTimeUse: storedSecret.IsOneTimeUse,
			Name: 	      storedSecret.Name,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
