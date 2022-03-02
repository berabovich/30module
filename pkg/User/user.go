package user

import (
	"fmt"
	"strings"
)

var Users []User

type User struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
}

func (u *User) FriendsToString() string {
	var friend []string
	var name string
	for _, fr := range u.Friends {
		for _, user := range Users {
			if user.Id == fr {
				name = user.Name
			}
		}
		friend = append(friend, name)
	}

	friends := strings.Join(friend, ", ")
	return friends
}

func (u *User) ToString() string {
	return fmt.Sprintf("id %s name %s age %s friends %s", u.Id, u.Name, u.Age, u.FriendsToString())
}
