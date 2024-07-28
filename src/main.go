package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/model"
)

var db = model.Database{
	Username: "grkmkly35",
	Server:   "",
	Database: "BookServices",
	Ctx:      context.TODO(),
}

func main() {
	godotenv.Load()

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connected")

	err = db.SetCollection("books")
	if err != nil {
		log.Fatal(err)
	}
	db.PrintDatabase()
	ID, _ := primitive.ObjectIDFromHex("66a631c024eb8e6658839297")
	uBook := model.Book{
		ObjectID: ID,
	}
	db.UpdateElementbyID("Notr Dame'nin Kamburu", &uBook)
	db.PrintDatabase()

}
