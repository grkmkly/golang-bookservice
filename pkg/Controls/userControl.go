package controls

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"main.go/model"
)

func GetHashedPassword(password string) (string, error) {
	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", errors.New("notGenerated")
	}
	return string(newHashPassword), nil
}

func CheckPassword(password string, newHashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(newHashPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
func HaveUsernameDB(db *model.Database, username string) (bool, error) {
	users, err := db.GetAllElements(model.User{})
	if err != nil {
		return false, errors.New("NotGetUser")
	}
	for _, v := range users {
		switch x := v.(type) {
		case model.User:
			if x.Username == username {
				return false, nil
			}
		}
	}
	return true, nil
}
