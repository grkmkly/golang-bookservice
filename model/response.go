package model

type Response struct {
	IsActive bool   `bson:"isActive"`
	Token    string `bson:"token"`
}
