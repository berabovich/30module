package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

//{
//	"name" : "some name",
//	"age" : 20,
//}

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}
type Friends struct {
	sourceId int `json:"source_id"`
	targetId int `json:"target_id"`
}

func (u *User) toString() string {
	return fmt.Sprintf("id %d name is %s and age is %d friends: %s \n", u.Id, u.Name, u.Age, u.Friends)
}

type service struct {
	store map[string]*User
}

func main() {

	mux := http.NewServeMux()
	srv := service{make(map[string]*User)}
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/make_friends", srv.MakeFriends)
	mux.HandleFunc("/get", srv.GetAll)
	http.ListenAndServe("localhost:8080", mux)
}
func (s *service) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		var u User
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		u.Id = 1
		for range s.store {
			u.Id++
		}

		s.store[strconv.Itoa(u.Id)] = &u

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User id " + strconv.Itoa(u.Id) + " status " + strconv.Itoa(http.StatusCreated)))

		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := ""
		for _, user := range s.store {
			response += user.toString()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))

		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
func (s *service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "MAKE_FRIENDS" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		var u User
		var f Friends
		if err := json.Unmarshal(content, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		//u.Friends[len(u.Friends)] = f.sourceId
		//u.Id = f.targetId
		//
		//s.store[strconv.Itoa(u.Id)] = &u
		for _, user := range s.store {
			if user.Id == f.targetId {
				user.Friends[len(user.Friends)] = f.sourceId
			}
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User name " + u.Name + " status " + strconv.Itoa(http.StatusCreated)))

		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
