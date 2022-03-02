package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
type UpdateUser struct {
	NewAge  string `json:"new_age"`
	NewName string `json:"new_name"`
}

var users []User

func (u *User) toString() string {
	var friend []string
	for _, fr := range u.Friends {
		f, _ := strconv.Atoi(fr)
		friend = append(friend, users[f-1].Name)
	}
	friends := strings.Join(friend, ", ")
	return fmt.Sprintf("id %s name %s age %s friends %s", u.Id, u.Name, u.Age, friends)
}

func main() {

	nr := chi.NewRouter()
	nr.MethodFunc("GET", "/users", getUsers)
	nr.MethodFunc("POST", "/create", createUser)
	nr.MethodFunc("GET", "/friends/{id}", getUserFriends)
	nr.MethodFunc("PUT", "/{id}", updateUserAge)
	nr.MethodFunc("DELETE", "/user", deleteUser)
	nr.MethodFunc("POST", "/make_friends", makeFriends)

	log.Fatal(http.ListenAndServe(":8080", nr))
}

func makeFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var makeFriend MakeFriends
	_ = json.NewDecoder(r.Body).Decode(&makeFriend)
	var name1 string
	var name2 string
	for index, u := range users {
		if u.Id == makeFriend.TargetId {
			users[index].Friends = append(users[index].Friends, makeFriend.SourceId)
			name1 = u.Name
		}
		if u.Id == makeFriend.SourceId {
			users[index].Friends = append(users[index].Friends, makeFriend.TargetId)
			name2 = u.Name
		}
	}

	w.Write([]byte("User " + name1 + " and User " + name2 + " now friends! Status: " + strconv.Itoa(http.StatusOK)))

}
func getUsers(w http.ResponseWriter, r *http.Request) {

	var response string
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
	for i, u := range users {
		if u.Id != strconv.Itoa(i+1) {
			user.Id = strconv.Itoa(i + 1)
			break
		}
		if u.Id == user.Id {
			id, _ := strconv.Atoi(user.Id)
			user.Id = strconv.Itoa(id + 1)
		}
	}
	users = append(users, user)
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})
	w.Write([]byte("User ID: " + user.Id + " Status:" + strconv.Itoa(http.StatusCreated)))
}
func getUserFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := chi.URLParam(r, "id")
	for _, item := range users {
		if item.Id == params {
			var friend []string
			for _, fr := range item.Friends {
				f, _ := strconv.Atoi(fr)
				friend = append(friend, users[f-1].Name)
			}
			friends := strings.Join(friend, ", ")
			w.Write([]byte("User: " + item.Name + " Friends: " + friends))
			break
		}
	}
}
func updateUserAge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updateUser UpdateUser
	_ = json.NewDecoder(r.Body).Decode(&updateUser)
	params := chi.URLParam(r, "id")
	for index, item := range users {
		if item.Id == params {
			users[index].Age = updateUser.NewAge
			w.Write([]byte("User " + item.Name + ". Age update successful! Status:" + strconv.Itoa(http.StatusOK)))
			return
		}
	}
	json.NewEncoder(w).Encode(users)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var makeFriend MakeFriends
	_ = json.NewDecoder(r.Body).Decode(&makeFriend)

	for i, u := range users {
		for j, f := range u.Friends {
			if f == makeFriend.TargetId {
				users[i].Friends = append(u.Friends[:j], u.Friends[j+1:]...)
			}
		}
	}
	for index, u := range users {
		if u.Id == makeFriend.TargetId {
			users = append(users[:index], users[index+1:]...)
			w.Write([]byte(u.Name + " was delete"))
			break
		}
	}

}
