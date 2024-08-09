package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Usersbook struct {
	UserId primitive.ObjectID `bson:"_userid"`
	BookId primitive.ObjectID `bson:"_bookid"`
}
