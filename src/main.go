package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

func createBook(w http.ResponseWriter, r *http.Request) {
	var b model.Book
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	b.ObjectID = primitive.NewObjectID()
	err = json.Unmarshal(body, &b)

	if err != nil {
		log.Fatal(err)
	}

	err = db.InsertDocument(&b)
	if err != nil {
		log.Fatal(err)
	}

}
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	bookId := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	books, err := db.GetAllElements()
	if err != nil {
		log.Fatal(err)
	}
	var updateBook model.Book

	err = json.Unmarshal(body, &updateBook)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range books {
		if v.ObjectID.Hex() == bookId {

			err := db.UpdateElementbyID(&updateBook, v)
			if err != nil {
				log.Fatal(err)
			}
			updateBook.ObjectID = v.ObjectID

			json, err := json.Marshal(&updateBook)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, string(json)+"%v", "updated")
			break
		}
	}

}
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := db.GetAllElements()
	if err != nil {
		log.Fatal(err)
	}
	json, err := json.Marshal(&books)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(json))
}

func getBook(w http.ResponseWriter, r *http.Request) {
	isHave := false
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	books, err := db.GetAllElements()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range books {
		if v.ObjectID.Hex() == vars["id"] {
			book, err := db.FindOneElementByID(v)
			if err != nil {
				log.Fatal(err)
			}
			json, err := json.Marshal(&book)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, string(json))
			isHave = true
			break
		}
	}
	if !isHave {
		fmt.Fprint(w, "NotFound")
	}
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	bookId := vars["id"]

	books, err := db.GetAllElements()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range books {
		if v.ObjectID.Hex() == bookId {
			err := db.DeleteElementByID(v)
			if err != nil {
				log.Fatal(err)
			}
			json, err := json.Marshal(v)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, string(json)+"%v", "deleted")
			break
		}
	}
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

	r := mux.NewRouter()
	r.HandleFunc("/createbook", createBook).Methods("POST")
	r.HandleFunc("/getbooks", getBooks).Methods("GET")
	r.HandleFunc("/getbook/{id}", getBook).Methods("GET")
	r.HandleFunc("/updatebook/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/deletebook/{id}", deleteBook).Methods("DELETE")

	http.ListenAndServe(":"+os.Getenv("PORT"), r)

}
