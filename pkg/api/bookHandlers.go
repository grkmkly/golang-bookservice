package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/model"
)

func createBook(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

}

func updateBook(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bookId := vars["id"]

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		books, err := db.GetAllElements(model.Book{})
		if err != nil {
			log.Fatal(err)
		}
		var updateBook model.Book

		err = json.Unmarshal(body, &updateBook)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range books {
			switch x := v.(type) {
			case model.Book:
				if x.ObjectID.Hex() == bookId {

					err := db.UpdateElementbyID(&updateBook, &x)
					if err != nil {
						log.Fatal(err)
					}
					updateBook.ObjectID = x.ObjectID

					json, err := json.Marshal(&updateBook)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Fprintf(w, string(json)+"%v", "updated")
					break
				}
			}
		}
	}
}

func getBooks(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		books, err := db.GetAllElements(model.Book{})
		if err != nil {
			log.Fatal(err)
		}
		json, err := json.Marshal(&books)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, string(json))
	}

}

func getBook(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isHave := false
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		books, err := db.GetAllElements(model.Book{})
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range books {
			switch x := v.(type) {
			case model.Book:
				if x.ObjectID.Hex() == vars["id"] {
					book, err := db.FindOneElementByID(&x)
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
		}
		if !isHave {
			fmt.Fprint(w, "NotFound")
		}
	}

}

func deleteBook(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bookId := vars["id"]

		books, err := db.GetAllElements(model.Book{})
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range books {
			switch x := v.(type) {
			case model.Book:
				if x.ObjectID.Hex() == bookId {
					err := db.DeleteElementByID(&x)
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
	}
}

func RoutesBook(r *mux.Router, db *model.Database) {
	r.HandleFunc("/createbook", createBook(db)).Methods("POST")
	r.HandleFunc("/getbooks", getBooks(db)).Methods("GET")
	r.HandleFunc("/getbook/{id}", getBook(db)).Methods("GET")
	r.HandleFunc("/updatebook/{id}", updateBook(db)).Methods("PUT")
	r.HandleFunc("/deletebook/{id}", deleteBook(db)).Methods("DELETE")
}
