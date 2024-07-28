package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ObjectID primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Author   string             `bson:"author"`
	Pages    string             `bson:"pages"`
	Topic    string             `bson:"topic"`
}
