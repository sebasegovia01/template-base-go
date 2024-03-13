package repositories

import (
	"context"
	"put-otp-go/src/models"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection_name = "otps"

var validate *validator.Validate

type OTPRepository struct {
	Database *mongo.Database
}

type IOTPRepository interface {
	Update(id string, otp models.Otp) (*mongo.UpdateResult, error)
}

func NewOTPRepository(db *mongo.Database) IOTPRepository {
	validate = validator.New()
	return &OTPRepository{Database: db}
}

// Update
func (br *OTPRepository) Update(id string, otp models.Otp) (*mongo.UpdateResult, error) {
	if err := validate.Struct(otp); err != nil {
		return nil, err
	}

	collection := br.Database.Collection(collection_name)

	// Change string id to primitive id
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": otp}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	return result, err
}
