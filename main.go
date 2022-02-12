package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const MIN_NAME_LENGTH = 5
const MIN_PASSWORD_LENGTH = 8

type UserBase struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserRegistration struct {
	UserBase
	Password string `json:"password"`
}

type UserSaved struct {
	UserBase
	HashedPassword string `json:"-"`
}

func (userData UserRegistration) verifyData() bool {
	if len(userData.Name) < MIN_NAME_LENGTH {
		return false
	}
	if len(userData.Password) < MIN_PASSWORD_LENGTH {
		return false
	}
	if len(userData.Email) < 4 {
		return false
	}

	return true
}

func (userData UserRegistration) getHashedPassword() (hashString string, err error) {
	//return userData.Password
	hashedBytes, errorHashing := bcrypt.GenerateFromPassword([]byte(userData.Password), 10)
	if errorHashing != nil {
		hashString = ""
		err = errorHashing
		return
	}

	hashString = string(hashedBytes)
	err = nil
	return
}

func (userData UserRegistration) createUserSaved() (savedUser UserSaved, err error) {
	fmt.Println("hit")
	if !userData.verifyData() {
		err = errors.New("User data is invalid")
		return
	}

	hashedPassword, errHashing := userData.getHashedPassword()
	if errHashing != nil {
		err = errHashing
		return
	}
	savedUser = UserSaved{
		UserBase{
			userData.Name,
			userData.Email,
		},
		hashedPassword,
	}
	err = nil
	return
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	var userData UserRegistration
	json.NewDecoder(r.Body).Decode(&userData)
	savedUser, err := userData.createUserSaved()
	if err != nil {
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{"Failed to create user"})
	} else {
		json.NewEncoder(w).Encode(&savedUser)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", handleRegistration).Methods("POST")
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", router)
}
