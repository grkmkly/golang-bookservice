package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"main.go/model"
)

func getHashedPassword(password string) (string, error) {
	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", errors.New("notGenerated")
	}
	return string(newHashPassword), nil
}

// func checkPassword(password string, newHashPassword string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(newHashPassword), []byte(password))
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

func registerUser(db *model.Database) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var u model.User

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Fatal(err)
		}

		u.ObjectID = primitive.NewObjectID()
		newPassword, err := getHashedPassword(u.Password)
		if err != nil {
			log.Fatal(err)
		}
		u.Password = newPassword

		for index := range u.Books {
			u.Books[index].ObjectID = primitive.NewObjectID()
		}
		err = db.InsertDocument(&u)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func RoutesUser(r *mux.Router, db *model.Database) {
	r.HandleFunc("/registeruser", registerUser(db)).Methods("POST")
	// r.HandleFunc("/getbooks", getBooks(&db)).Methods("GET")
	// r.HandleFunc("/getbook/{id}", getBook(&db)).Methods("GET")
	// r.HandleFunc("/updatebook/{id}", updateBook(&db)).Methods("PUT")
	// r.HandleFunc("/deletebook/{id}", deleteBook(&db)).Methods("DELETE")
}
