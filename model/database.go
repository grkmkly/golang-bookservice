package model

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Type Database
type Database struct {
	Username   string
	Server     string
	Database   string
	Collection *mongo.Collection
	Client     *mongo.Client
	Ctx        context.Context
}

// Connect Database
func (db *Database) Connect() error {
	uri := os.Getenv("mongoUri")
	if uri == "" {
		return errors.New("UriNotFound")
	}
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(db.Ctx, clientOptions)

	if err != nil {
		return err
	}
	db.Client = client
	err = client.Ping(db.Ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// Set Collection
func (db *Database) SetCollection(coll string) error {
	if coll == "" {
		return errors.New("CollectionisNull")
	}
	collection := db.Client.Database(db.Database).Collection(coll)
	db.Collection = collection
	return nil
}

// Insert Document
func (db *Database) InsertDocument(mBook ...*Book) error {
	if len(mBook) == 0 {
		return errors.New("UnexpectedLength")
	}
	for _, value := range mBook {
		_, err := db.Collection.InsertOne(db.Ctx, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) FindOneElementByID(mbook *Book) (*Book, error) {
	filter := bson.D{
		{Key: "_id", Value: mbook.ObjectID},
	}
	result := db.Collection.FindOne(db.Ctx, filter)
	if result == nil {
		return nil, errors.New("ResultNotFound")
	}

	err := result.Decode(mbook)
	if err != nil {
		return nil, err
	}

	return mbook, nil
}

// Get Element
func (db *Database) GetFullElement() ([]*Book, error) {
	filter := bson.D{}
	cursor, err := db.Collection.Find(db.Ctx, filter)
	if err != nil {
		return nil, err
	}
	var books []*Book
	for cursor.Next(db.Ctx) {
		var result Book
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		books = append(books, &result)
	}
	return books, nil
}
