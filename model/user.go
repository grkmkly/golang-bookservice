package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ObjectID primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Books    []Book             `bson:"books"`
}
