package model

type Response struct {
	IsActive bool   `bson:"isActive"`
	Token    string `bson:"token"`
	ID       string `bson:"_id"`
	Username string `bson:"username"`
}
