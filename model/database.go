package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (db *Database) InsertDocument(item interface{}) error {
	typeInterface := db.ControlItemTypeandSet(item)
	switch typeInterface {
	case reflect.TypeOf(Book{}):

		_, err := db.Collection.InsertOne(db.Ctx, item)
		if err != nil {
			return err
		}

	case reflect.TypeOf(User{}):
		_, err := db.Collection.InsertOne(db.Ctx, item)
		if err != nil {
			return err
		}
	default:
		return errors.New("TypeNotFound")
	}
	return nil
}

func (db *Database) ControlItemTypeandSet(item interface{}) interface{} {

	ft := reflect.TypeOf(item)
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	switch ft {
	case reflect.TypeOf(Book{}):
		db.SetCollection("books")
	case reflect.TypeOf(User{}):
		db.SetCollection("users")
	}
	return ft
}

// Finding element by ID
func (db *Database) FindOneElementByID(item interface{}) (interface{}, error) {
	//only set
	db.ControlItemTypeandSet(item)

	if x, ok := item.(Book); ok {
		fmt.Print("Selam")
		var book Book
		filter := bson.D{
			{Key: "_id", Value: x.ObjectID},
		}
		result := db.Collection.FindOne(db.Ctx, filter)
		if result == nil {
			return nil, errors.New("ResultNotFound")
		}
		err := result.Decode(&book)
		if err != nil {
			return nil, err
		}
		return book, nil
	} else if x, ok := item.(User); ok {
		var user User
		filter := bson.D{
			{Key: "_id", Value: x.ObjectID},
		}
		result := db.Collection.FindOne(db.Ctx, filter)
		if result == nil {
			return nil, errors.New("ResultNotFound")
		}
		err := result.Decode(&user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, errors.New("ItemNotFound")
}

// Get Element
func (db *Database) GetAllElements(item interface{}) ([]interface{}, error) {

	typeInterface := db.ControlItemTypeandSet(item)

	filter := bson.D{}
	cursor, err := db.Collection.Find(db.Ctx, filter)
	if err != nil {
		return nil, err
	}

	switch typeInterface {
	case reflect.TypeOf(Book{}):
		var bookItems []interface{}
		for cursor.Next(db.Ctx) {
			var result Book
			err = cursor.Decode(&result)
			if err != nil {
				return nil, err
			}
			bookItems = append(bookItems, result)
		}
		return bookItems, nil
	case reflect.TypeOf(User{}):
		var userItems []interface{}
		for cursor.Next(db.Ctx) {
			var result User
			err = cursor.Decode(&result)
			if err != nil {
				return nil, err
			}
			userItems = append(userItems, result)
		}
		return userItems, nil
	default:
		return nil, errors.New("TypeProblem")
	}
}

// Delete Element
func (db *Database) DeleteElementByID(item interface{}) error {
	// only set
	db.ControlItemTypeandSet(item)

	switch x := item.(type) {
	case Book:
		filter := bson.D{
			{Key: "_id", Value: x.ObjectID},
		}
		_, err := db.Collection.DeleteOne(db.Ctx, filter)
		if err != nil {
			return err
		}
	case User:
		filter := bson.D{
			{Key: "_id", Value: x.ObjectID},
		}
		_, err := db.Collection.DeleteOne(db.Ctx, filter)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update Element
func (db *Database) UpdateElementbyID(updateBook *Book, filterBook *Book) error {
	if (filterBook == nil) || (updateBook == nil) {
		return errors.New("BookIsNull")
	}
	filter := bson.D{
		{Key: "_id", Value: filterBook.ObjectID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: updateBook.Name},
			{Key: "author", Value: updateBook.Author},
			{Key: "pages", Value: updateBook.Pages},
			{Key: "topic", Value: updateBook.Topic},
		}},
	}
	_, err := db.Collection.UpdateOne(db.Ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) AddbookUser(user User, book Book) error {
	db.ControlItemTypeandSet(user)
	u, err := db.FindOneElementByID(user)
	if err != nil {
		return err
	}
	var filterUser primitive.D
	switch x := u.(type) {
	case User:
		filterUser = bson.D{
			{Key: "_id", Value: x.ObjectID},
		}
	}
	switch x := u.(type) {
	case User:
		x.Books = append(x.Books, book)
		updateUser := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "Books", Value: x.Books},
			}},
		}
		_, err := db.Collection.UpdateOne(db.Ctx, filterUser, updateUser)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// Prints the Database
func (db *Database) PrintDatabase(item interface{}) error {
	items, err := db.GetAllElements(item)
	if err != nil {
		return errors.New("NotGetElement")
	}
	for _, v := range items {
		switch x := v.(type) {
		case Book:
			pattern := fmt.Sprintf("Id : %v\nName : %v\nAuthor : %v\nPages : %v\nTopic : %v\n", x.ObjectID, x.Name, x.Author, x.Pages, x.Topic)
			fmt.Println(pattern)
		case User:
			pattern := fmt.Sprintf("Id : %v\nUserame : %v\nBooks : %v\n", x.ObjectID, x.Username, x.Books)
			fmt.Println(pattern)
		default:
			fmt.Println("NotFoundType")
		}
	}
	return nil
}
