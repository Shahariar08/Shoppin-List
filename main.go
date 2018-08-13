package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var ID_generator int //a generator variable which genereates unique id for an item
type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Price    int    `json:"price"`
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
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(item)
	if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!an Error Occured"})
			return 
	}

}

func additem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var curitem Item
	err := json.NewDecoder(r.Body).Decode(&curitem)
	if err == nil {
		if curitem.Name == "" || curitem.Quantity == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!insufficiet information"})
			return
		}
		ID_generator++
		curitem.ID = ID_generator
		item = append(item, curitem)
		err2 := json.NewEncoder(w).Encode(curitem)
		if err2 != nil{
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!an Error Occured"})
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Ok: 0, Message: "An error Occured"})
		return
	}
}
func deleteitem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//parameters have been handled
	gotId, err := strconv.Atoi(params["id"])

	if err == nil {
		for index, i := range item {
			if i.ID == gotId {
				item = append(item[:index], item[index+1:]...)
				//found the desired index , removed and then appended the rest of the items items
				json.NewEncoder(w).Encode(item)
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
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	gotId, err := strconv.Atoi(params["id"])
	if err == nil {
		for index, i := range item {
			if i.ID == gotId {
				item = append(item[:index], item[index+1:]...)
				//basically deletion

				var curitem Item
				curitem = i
				err2 := json.NewDecoder(r.Body).Decode(&curitem)
				if err2 == nil {
					item = append(item, curitem)
					//updated item has been added

					err3 := json.NewEncoder(w).Encode(item)
					if err3 != nil {
						w.WriteHeader(http.StatusBadRequest)
		                json.NewEncoder(w).Encode(Response{Ok: 0, Message: "An error Occured"})
		                return
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
		return ;
	}
}

func main() {
	ID_generator = 1
	//sample items added
	//item = append(item, Item{ID: 1, Name: "Fish", Quantity: "5kg", Price: 350})
	//item = append(item, Item{ID: 2, Name: "Rice", Quantity: "2kg", Price: 120})
	// item = append(item , Item{ID:3 , Name:"fuck", Quantity:5 , Price : 350})
	m := mux.NewRouter()
	//router
	//function handler
	m.HandleFunc("/", homepage).Methods("GET")
	m.HandleFunc("/showitem", showitem).Methods("GET")
	m.HandleFunc("/additem", additem).Methods("POST")
	m.HandleFunc("/deleteitem/{id}", deleteitem).Methods("DELETE")
	m.HandleFunc("/updateitem/{id}", updateitem).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8800", m))
}
