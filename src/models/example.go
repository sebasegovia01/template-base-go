package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Book  model example
type Book struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title" validate:"required"`
	Author string             `bson:"author" validate:"required"`
}
