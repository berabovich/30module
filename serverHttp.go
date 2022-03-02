package main

import (
	"30module/pkg/HttpReqests"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

func main() {

	nr := chi.NewRouter()
	nr.MethodFunc("GET", "/users", HttpReqests.HttpGetUsers)
	nr.MethodFunc("POST", "/create", HttpReqests.HttpCreateUser)
	nr.MethodFunc("GET", "/friends/{id}", HttpReqests.HttpGetUserFriends)
	nr.MethodFunc("PUT", "/{id}", HttpReqests.HttpUpdateUserAge)
	nr.MethodFunc("DELETE", "/user", HttpReqests.HttpDeleteUser)
	nr.MethodFunc("POST", "/make_friends", HttpReqests.HttpMakeFriends)

	log.Fatal(http.ListenAndServe(":8080", nr))
}
