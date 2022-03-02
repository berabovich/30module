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

func (u *User) friendsToString() string {
	var friend []string
	var name string
	for _, fr := range u.Friends {
		for _, user := range users {
			if user.Id == fr {
				name = user.Name
			}
		}
		friend = append(friend, name)
	}

	friends := strings.Join(friend, ", ")
	return friends
}

func (u *User) toString() string {
	return fmt.Sprintf("id %s name %s age %s friends %s", u.Id, u.Name, u.Age, u.friendsToString())
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
	_, err := strconv.Atoi(makeFriend.TargetId)
	if err != nil {
		_, err := w.Write([]byte("ID must be int type"))
		if err != nil {
			return
		}
		return
	}
	_, err = strconv.Atoi(makeFriend.SourceId)
	if err != nil {
		_, err := w.Write([]byte("ID must be int type"))
		if err != nil {
			return
		}
		return
	}
	for _, u := range users {
		if u.Id == makeFriend.TargetId {
			name1 = u.Name
		}
		if u.Id == makeFriend.SourceId {
			name2 = u.Name
		}
	}
	if name1 == "" || name2 == "" {
		_, err := w.Write([]byte("Users not found"))
		if err != nil {
			return
		}
		return
	}
	for index, u := range users {
		if u.Id == makeFriend.TargetId {
			users[index].Friends = append(users[index].Friends, makeFriend.SourceId)
		}
		if u.Id == makeFriend.SourceId {
			users[index].Friends = append(users[index].Friends, makeFriend.TargetId)
		}
	}

	_, err = w.Write([]byte("User " + name1 + " and User " + name2 + " now friends! Status: " + strconv.Itoa(http.StatusOK)))
	if err != nil {
		return
	}
}
func getUsers(w http.ResponseWriter, _ *http.Request) {

	var response string
	for _, user := range users {
		response += user.toString() + "\n"
	}
	_, err := w.Write([]byte(response))
	if err != nil {
		return
	}

}
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	_, err := strconv.Atoi(user.Age)
	if err != nil {
		_, err := w.Write([]byte("Age must be int type"))
		if err != nil {
			return
		}
		return
	}
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
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("User ID: " + user.Id + " Status:" + strconv.Itoa(http.StatusCreated)))
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
func getUserFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := chi.URLParam(r, "id")
	for _, u := range users {
		if u.Id == params {
			_, err := w.Write([]byte("User: " + u.Name + " Friends: " + u.friendsToString()))
			if err != nil {
				return
			}
			return
		}
	}
	_, err := w.Write([]byte("User not find"))
	if err != nil {
		return
	}
}
func updateUserAge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updateUser UpdateUser
	_ = json.NewDecoder(r.Body).Decode(&updateUser)
	params := chi.URLParam(r, "id")
	_, err := strconv.Atoi(updateUser.NewAge)
	if err != nil {
		_, err := w.Write([]byte("Age must be int type"))
		if err != nil {
			return
		}
		return
	}
	for index, item := range users {
		if item.Id == params {
			users[index].Age = updateUser.NewAge
			_, err := w.Write([]byte("User " + item.Name + ". Age update successful! Status: " + strconv.Itoa(http.StatusOK)))
			if err != nil {
				return
			}
			return
		}
	}
	_, err = w.Write([]byte("User not found"))
	if err != nil {
		return
	}
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
			_, err := w.Write([]byte(u.Name + " was delete. Status: " + strconv.Itoa(http.StatusOK)))
			if err != nil {
				return
			}
			return
		}
	}
	_, err := w.Write([]byte("User not found"))
	if err != nil {
		return
	}
}
