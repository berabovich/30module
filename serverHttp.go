package main

import (
	"30module/pkg/httpReqests"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

func main() {

	nr := chi.NewRouter()
	nr.MethodFunc("GET", "/users", httpRequests.HttpGetUsers)
	nr.MethodFunc("POST", "/create", httpRequests.HttpCreateUser)
	nr.MethodFunc("GET", "/friends/{id}", httpRequests.HttpGetUserFriends)
	nr.MethodFunc("PUT", "/{id}", httpRequests.HttpUpdateUserAge)
	nr.MethodFunc("DELETE", "/user", httpRequests.HttpDeleteUser)
	nr.MethodFunc("POST", "/make_friends", httpRequests.HttpMakeFriends)

	log.Fatal(http.ListenAndServe(":8080", nr))
}
