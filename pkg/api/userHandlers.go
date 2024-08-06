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
	controls "main.go/pkg/Controls"
)

func registerUser(db *model.Database) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var u model.User
		var resp model.Response
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Fatal(err)
		}
		isHave, err := controls.HaveUsernameDB(db, u.Username)
		if err != nil {
			log.Fatal(err)
		}
		switch isHave {
		case true:
			fmt.Println("UsernameAvailable")
		case false:
			fmt.Println("UsernameNotAvailable")
			resp.IsActive = false
			jsonResponse, err := json.Marshal(&resp)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, string(jsonResponse))
			return
		}
		u.ObjectID = primitive.NewObjectID()
		newPassword, err := controls.GetHashedPassword(u.Password)
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
		resp.IsActive = true
		jsonResponse, err := json.Marshal(&resp)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, string(jsonResponse))
	}
}
func signinUser(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		isCheck := false
		var resp model.Response
		var u model.User
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Fatal(err)
		}
		users, err := db.GetAllElements(model.User{})
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range users {
			switch x := v.(type) {
			case model.User:
				if (x.Username == u.Username) && controls.CheckPassword(u.Password, x.Password) {
					isCheck = true
					break
				}
			default:
				fmt.Print("NotUser")
				return
			}
		}
		if isCheck {
			resp.IsActive = true
			jsonResponse, err := json.Marshal(&resp)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, string(jsonResponse))
		} else {
			resp.IsActive = false
			jsonResponse, err := json.Marshal(&resp)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, string(jsonResponse))
		}
	}
}

func RoutesUser(r *mux.Router, db *model.Database) {
	r.HandleFunc("/registeruser", registerUser(db)).Methods("POST")
	r.HandleFunc("/signinuser", signinUser(db)).Methods("POST")
	//r.HandleFunc("/getbook/{id}", getBook(db)).Methods("GET")
	//r.HandleFunc("/updatebook/{id}", updateBook(db)).Methods("PUT")
	//r.HandleFunc("/deletebook/{id}", deleteBook(db)).Methods("DELETE")
}
