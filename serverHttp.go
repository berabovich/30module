package main

import (
	"30module/pkg/HttpReqests"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

func main() {

	nr := chi.NewRouter()
	nr.MethodFunc("GET", "/users", httpReqests.HttpGetUsers)
	nr.MethodFunc("POST", "/create", httpReqests.HttpCreateUser)
	nr.MethodFunc("GET", "/friends/{id}", httpReqests.HttpGetUserFriends)
	nr.MethodFunc("PUT", "/{id}", httpReqests.HttpUpdateUserAge)
	nr.MethodFunc("DELETE", "/user", httpReqests.HttpDeleteUser)
	nr.MethodFunc("POST", "/make_friends", httpReqests.HttpMakeFriends)

	log.Fatal(http.ListenAndServe(":8080", nr))
}
