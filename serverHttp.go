package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
}
type MakeFriends struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
}

var makeFriend []MakeFriends

var users []User

func (u *User) toString() string {
	return fmt.Sprintf("id %s name %s age %s friends %s", u.Id, u.Name, u.Age, u.Friends)
}

//func (m *MakeFriends) toString() string {
//	return fmt.Sprintf("target %s source %s", m.TargetId, m.SourceId)
//}

func main() {

	nr := chi.NewRouter()
	nr.MethodFunc("GET", "/users", getUsers)
	nr.MethodFunc("POST", "/create", createUser)
	nr.MethodFunc("GET", "/users/{id}", getUser)
	nr.MethodFunc("PUT", "/{id}", updateUserAge)
	nr.MethodFunc("DELETE", "/{id}", deleteUser)
	nr.MethodFunc("POST", "/make_friends", makeFriends)

	log.Fatal(http.ListenAndServe(":8080", nr))
}

func makeFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user1 string
	var user2 string
	var makeFriend MakeFriends
	_ = json.NewDecoder(r.Body).Decode(&makeFriend)
	for index, item := range users {
		if item.Id == makeFriend.TargetId {
			users = append(users[:index], users[index+1:]...)
			var user User
			user.Name = item.Name
			user.Friends[index] = makeFriend.SourceId
			user.Id = item.Id
			user1 = user.Name
			users = append(users, user)
			return
		}
		//if item.Id == makeFriend.SourceId {
		//	users = append(users[:index], users[index+1:]...)
		//	var user User
		//	_ = json.NewDecoder(r.Body).Decode(&user)
		//	user.Name = item.Name
		//	user.Friends[index] = makeFriend.TargetId
		//	user.Id = item.Id
		//	user2 = user.Name
		//	users = append(users, user)
		//	return
		//}
	}

	w.Write([]byte("User " + user1 + " and user " + user2 + " now friends! Status: " + strconv.Itoa(http.StatusOK)))

}
func getUsers(w http.ResponseWriter, r *http.Request) {

	response := ""
	for _, user := range users {
		response += user.toString() + "\n"
	}
	w.Write([]byte(response))

}
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.Id = strconv.Itoa(len(users) + 1)
	users = append(users, user)
	w.Write([]byte("Used ID: " + user.Id + " Status:" + strconv.Itoa(http.StatusCreated)))
}
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := chi.URLParam(r, "id")
	for _, item := range users {
		if item.Id == params {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}
func updateUserAge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := chi.URLParam(r, "id")
	for index, item := range users {
		if item.Id == params {
			users = append(users[:index], users[index+1:]...)
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)
			user.Name = item.Name
			user.Friends = item.Friends
			user.Id = params
			users = append(users, user)
			w.Write([]byte("User age update successful! Status:" + strconv.Itoa(http.StatusOK)))
			return
		}
	}
	json.NewEncoder(w).Encode(users)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := chi.URLParam(r, "id")
	for index, item := range users {
		if item.Id == params {
			users = append(users[:index], users[index+1:]...)
			w.Write([]byte(item.Name + " was delete"))
			break
		}
	}
}
