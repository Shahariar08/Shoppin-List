package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

/************************************Authorization starts here********************************/
type User struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Ok          bool
	Message     string
	Information string
}

//A map that stores user informations globally
var userList = make(map[string]User)

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("User")
	if err == nil {
		_, flag := userList[cookie.Value] //checking cookie value .if it matches with the current user name it returns true
		return flag
	}
	return false
}

func regUser(w http.ResponseWriter, r *http.Request) {
	//Registration can not be performed being logged in
	//checking if logged in
	//Post method was used
	if isLoggedIn(w, r) == true {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserResponse{false, "Please logout to register!", ""})
		return
	}
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserResponse{false, "Invalid request!", ""})
	} else {
		if _, found := userList[newUser.UserName]; found == true {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(UserResponse{false, "User already exists", newUser.UserName})
		} else if newUser.UserName == "" || newUser.Password == "" || newUser.Name == "" {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(UserResponse{false, "Invalid user info", ""})
		} else {
			userList[newUser.UserName] = newUser
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(UserResponse{true, "Registered new user", newUser.UserName})
		}
	}
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	//Post method was used
	userName, password, flag := r.BasicAuth()
	if isLoggedIn(w, r) == true {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(UserResponse{false, "Already logged in", ""})
		return
	}
	if flag == false {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserResponse{false, "Invalid request!", ""})
	} else {
		val, found := userList[userName]
		if found == true && val.Password == password {
			cookie := http.Cookie{Name: "User", Value: userName, Path: "/"}
			//setting cookie value
			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(UserResponse{true, "Successfully logged in", userName})
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(UserResponse{false, "Invalid username or password", userName})
		}
	}
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("User")
	if err == nil {
		cookie := http.Cookie{Name: "User", Value: "", Path: "/", Expires: time.Now()}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(UserResponse{true, "Logged out", ""})
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(UserResponse{false, "No active user found", ""})
	}
}

/************************************Autorization ends here*********************/

/***********************************Shopping List starts here**************/

var ID_generator int //a generator variable which genereates unique id for an item
type Item struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	PricePerUnit int    `json:"pricePerUnit"`
	TotalPrice   int    `json:"totalPrice"`
}
type Response struct {
	Ok      int    `json : "ok"`
	Message string `json : "message"`
}

var item []Item

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Shopping List RESTful api server")
}

func showitem(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(w, r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{0, "Please login first"})
		return
	}

	err := json.NewEncoder(w).Encode(item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!an Error Occured"})
		return
	}
}

func additem(w http.ResponseWriter, r *http.Request) {

	if isLoggedIn(w, r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{0, "Please login first"})
		return
	}

	var curitem Item
	err := json.NewDecoder(r.Body).Decode(&curitem)
	if err == nil {
		if curitem.Name == "" || curitem.Quantity <= 0 || curitem.PricePerUnit <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!insufficiet or invalid information"})
			return
		}
		for index, i := range item {
			if strings.ToLower(i.Name) == strings.ToLower(curitem.Name) { //check if an item already exists or not making all the characters to lower case
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!This item name already exists.Please try to update it with a put request"})
				return
			}
			item[index] = item[index]
		}

		ID_generator++
		curitem.ID = ID_generator
		curitem.TotalPrice = curitem.PricePerUnit * curitem.Quantity

		item = append(item, curitem)
		err2 := json.NewEncoder(w).Encode(curitem)
		if err2 != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!an Error Occured"})
		}

		json.NewEncoder(w).Encode(Response{Ok: 1, Message: "Item has been added Sucessfully!"})

	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Ok: 0, Message: "An error Occured"})
		return
	}
}
func deleteitem(w http.ResponseWriter, r *http.Request) {

	if isLoggedIn(w, r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{0, "Please login first"})
		return
	}

	params := mux.Vars(r)
	//parameters have been handled
	gotId, err := strconv.Atoi(params["id"])

	if err == nil {
		for index, i := range item {
			if i.ID == gotId {
				item = append(item[:index], item[index+1:]...)
				//found the desired index , removed and then appended the rest of the items items
				err2 := json.NewEncoder(w).Encode(item)
				if err2 == nil {
					json.NewEncoder(w).Encode(Response{Ok: 1, Message: "Item has been Deleted Sucessfully!"})
				} else {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(Response{Ok: 0, Message: "An error Occured"})
					return
				}
				return
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!Such ID does not exist in the list"})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Ok: 0, Message: "An error Occured"})
		return
	}

}
func updateitem(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(w, r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{0, "Please login first"})
		return
	}

	params := mux.Vars(r)
	gotId, err := strconv.Atoi(params["id"])
	if err == nil {
		for index, i := range item {
			if i.ID == gotId {

				var curitem Item
				err2 := json.NewDecoder(r.Body).Decode(&curitem)
				if err2 == nil {
					if curitem.Name == "" || curitem.Quantity <= 0 || curitem.PricePerUnit <= 0 {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!insufficiet or invalid information"})
						return
					}
					curitem.ID = gotId
					curitem.TotalPrice = curitem.PricePerUnit * curitem.Quantity
					item[index] = curitem //update the current item
					err3 := json.NewEncoder(w).Encode(item)
					if err3 != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(Response{Ok: 0, Message: "An error Occured"})
						return
					} else {
						json.NewEncoder(w).Encode(Response{Ok: 1, Message: "Item has been updated Sucessfully!"})
					}

				} else {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(Response{Ok: 0, Message: "An error Occured"})
				}
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!Such ID does not exist in the list"})

	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Ok: 0, Message: "An error Occured"})
		return
	}
}

/***********************************Shopping List ends here**************/

func main() {
	ID_generator = 0

	m := mux.NewRouter() //router
	m.HandleFunc("/", homepage).Methods("GET")
	m.HandleFunc("/showitem", showitem).Methods("GET")
	m.HandleFunc("/additem", additem).Methods("POST")
	m.HandleFunc("/deleteitem/{id}", deleteitem).Methods("DELETE")
	m.HandleFunc("/updateitem/{id}", updateitem).Methods("PUT")

	m.HandleFunc("/register", regUser).Methods("POST")
	m.HandleFunc("/login", loginUser).Methods("POST")
	m.HandleFunc("/logout", logoutUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":8800", m))
}
