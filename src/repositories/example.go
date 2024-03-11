package repositories

import (
	"context"
	"template-base-go/src/models"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection_name = "books"

var validate *validator.Validate

type ExampleRepository struct {
	Database *mongo.Database
}

type IExampleRepository interface {
	Create(book models.Book) (*mongo.InsertOneResult, error)
	Update(id string, book models.Book) (*mongo.UpdateResult, error)
	Get(id string) (*models.Book, error)
}

func NewExampleRepository(db *mongo.Database) IExampleRepository {
	validate = validator.New()
	return &ExampleRepository{Database: db}
}

// Create
func (br *ExampleRepository) Create(book models.Book) (*mongo.InsertOneResult, error) {

	if err := validate.Struct(book); err != nil {
		return nil, err
	}

	collection := br.Database.Collection(collection_name)
	result, err := collection.InsertOne(context.Background(), book)
	return result, err
}

// Update
func (br *ExampleRepository) Update(id string, book models.Book) (*mongo.UpdateResult, error) {
	if err := validate.Struct(book); err != nil {
		return nil, err
	}

	collection := br.Database.Collection(collection_name)

	// Change string id to primitive id
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": book}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	return result, err
}

// Get
func (br *ExampleRepository) Get(id string) (*models.Book, error) {
	collection := br.Database.Collection(collection_name)
	var book models.Book

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	err = collection.FindOne(context.Background(), filter).Decode(&book)
	return &book, err
}
