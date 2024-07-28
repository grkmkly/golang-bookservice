package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
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

	if err != nil {
		log.Fatal(err)
	}
}
