package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"main.go/model"
	"main.go/pkg/api"
)

var db = &model.Database{
	Username: "grkmkly35",
	Server:   "",
	Database: "BookServices",
	Ctx:      context.TODO(),
}

func deneme(item interface{}) {
	switch x := item.(type) {
	case model.Book:
		fmt.Print(x.Name)
	case model.User:
		fmt.Print(x.Username)
	}
}

func main() {

	godotenv.Load()
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connected")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	db.SetCollection("books")
	api.RoutesBook(r, db)
	api.RoutesUser(r, db)

	http.ListenAndServe(":"+os.Getenv("PORT"), r)

}
