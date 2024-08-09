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
		case false:
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
		resp.ID = u.ObjectID.Hex()
		resp.Username = u.Username
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
					resp.ID = x.ObjectID.Hex()
					resp.Username = x.Username
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

func getUsers(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		users, err := db.GetAllElements(model.User{})
		if err != nil {
			log.Fatal(err)
		}
		var username []model.User
		for _, v := range users {
			switch x := v.(type) {
			case model.User:
				models := model.User{
					ObjectID: x.ObjectID,
					Username: x.Username,
					Books:    x.Books,
				}
				username = append(username, models)
			}
		}
		jsonUser, _ := json.Marshal(&username)
		fmt.Fprint(w, string(jsonUser))
	}
}
func getUser(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		if vars["id"] == "" {
			fmt.Fprint(w, "NotFoundId")
			return
		}
		var err error
		//resp, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		var user model.User
		user.ObjectID, err = primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			log.Fatal(err)
		}

		// err = json.Unmarshal(resp, &user)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		findUser, err := db.FindOneElementByID(user)
		if err != nil {
			log.Fatal(err)
		}
		switch x := findUser.(type) {
		case model.User:
			var model = model.User{
				ObjectID: x.ObjectID,
				Username: x.Username,
				Books:    x.Books,
			}
			jsonResponse, err := json.Marshal(model)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, string(jsonResponse))
		}
	}
}
func addbookUser(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		isCheck := false
		// Read Body
		jsonRequest, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// Body to struct
		var usersbook model.Usersbook
		err = json.Unmarshal(jsonRequest, &usersbook)
		if err != nil {
			log.Fatal(err)
		}
		u := model.User{
			ObjectID: usersbook.UserId,
		}
		b := model.Book{
			ObjectID: usersbook.BookId,
		}
		book, err := db.FindOneElementByID(b)
		if err != nil {
			log.Fatal(err)
		}
		switch x := book.(type) {
		case model.Book:
			err = db.AddbookUser(u, x)
			if err != nil {
				log.Fatal(err)
			}
			isCheck = true
		}
		var resp model.Response
		resp.IsActive = isCheck
		resp.ID = u.Username
		jsonResponse, err := json.Marshal(resp)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, string(jsonResponse))
	}
}

func RoutesUser(r *mux.Router, db *model.Database) {
	r.HandleFunc("/registeruser", registerUser(db)).Methods("POST")
	r.HandleFunc("/signinuser", signinUser(db)).Methods("POST")
	r.HandleFunc("/getusers", getUsers(db)).Methods("GET")
	r.HandleFunc("/getuser/{id}", getUser(db)).Methods("GET")
	r.HandleFunc("/addbookuser", addbookUser(db)).Methods("POST")
}
