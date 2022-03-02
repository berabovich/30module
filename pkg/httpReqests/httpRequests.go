package httpRequests

import (
	"30module/pkg/makeFriends"
	"30module/pkg/updateUser"
	"30module/pkg/user"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"sort"
	"strconv"
)

func HttpMakeFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var makeFriend makeFriends.MakeFriends
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
	for _, u := range user.Users {
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
	for index, u := range user.Users {
		if u.Id == makeFriend.TargetId {
			user.Users[index].Friends = append(user.Users[index].Friends, makeFriend.SourceId)
		}
		if u.Id == makeFriend.SourceId {
			user.Users[index].Friends = append(user.Users[index].Friends, makeFriend.TargetId)
		}
	}

	_, err = w.Write([]byte("User " + name1 + " and User " + name2 + " now friends! Status: " + strconv.Itoa(http.StatusOK)))
	if err != nil {
		return
	}
}
func HttpGetUsers(w http.ResponseWriter, _ *http.Request) {

	var response string
	for _, u := range user.Users {
		response += u.ToString() + "\n"
	}
	_, err := w.Write([]byte(response))
	if err != nil {
		return
	}

}
func HttpCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u user.User
	_ = json.NewDecoder(r.Body).Decode(&u)
	_, err := strconv.Atoi(u.Age)
	if err != nil {
		_, err := w.Write([]byte("Age must be int type"))
		if err != nil {
			return
		}
		return
	}
	u.Id = strconv.Itoa(len(user.Users) + 1)
	for i, u := range user.Users {
		if u.Id != strconv.Itoa(i+1) {
			u.Id = strconv.Itoa(i + 1)
			break
		}
		if u.Id == u.Id {
			id, _ := strconv.Atoi(u.Id)
			u.Id = strconv.Itoa(id + 1)
		}
	}
	user.Users = append(user.Users, u)
	sort.SliceStable(user.Users, func(i, j int) bool {
		return user.Users[i].Id < user.Users[j].Id
	})
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("User ID: " + u.Id + " Status:" + strconv.Itoa(http.StatusCreated)))
	if err != nil {
		return
	}
}
func HttpGetUserFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := chi.URLParam(r, "id")
	for _, u := range user.Users {
		if u.Id == params {
			_, err := w.Write([]byte("User: " + u.Name + " Friends: " + u.FriendsToString()))
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
func HttpUpdateUserAge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updateAge updateUser.UpdateUser
	_ = json.NewDecoder(r.Body).Decode(&updateAge)
	params := chi.URLParam(r, "id")
	_, err := strconv.Atoi(updateAge.NewAge)

	if err != nil {
		_, err := w.Write([]byte("Age must be int type"))
		if err != nil {
			return
		}
		return
	}
	for index, item := range user.Users {
		if item.Id == params {
			user.Users[index].Age = updateAge.NewAge
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
func HttpDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var makeFriend makeFriends.MakeFriends
	_ = json.NewDecoder(r.Body).Decode(&makeFriend)

	for i, u := range user.Users {
		for j, f := range u.Friends {
			if f == makeFriend.TargetId {
				user.Users[i].Friends = append(u.Friends[:j], u.Friends[j+1:]...)
			}
		}
	}
	for index, u := range user.Users {
		if u.Id == makeFriend.TargetId {
			user.Users = append(user.Users[:index], user.Users[index+1:]...)
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
