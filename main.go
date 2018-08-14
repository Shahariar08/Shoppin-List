package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

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
	err := json.NewEncoder(w).Encode(item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!an Error Occured"})
		return
	}
}

func additem(w http.ResponseWriter, r *http.Request) {
	var curitem Item
	err := json.NewDecoder(r.Body).Decode(&curitem)
	if err == nil {
		if curitem.Name == "" || curitem.Quantity <= 0 || curitem.PricePerUnit <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!insufficiet or invalid information"})
			return
		}
		for index, i := range item {
			if strings.ToLower(i.Name) == strings.ToLower(curitem.Name) {
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
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	gotId, err := strconv.Atoi(params["id"])
	if err == nil {
		for index, i := range item {
			if i.ID == gotId {
				//item = append(item[:index], item[index+1:]...)
				//basically deletion

				var curitem Item
				err2 := json.NewDecoder(r.Body).Decode(&curitem)
				if err2 == nil {
					//item = append(item, curitem)
					//updated item has been added
					if curitem.Name == "" || curitem.Quantity <= 0 || curitem.PricePerUnit <= 0 {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(Response{Ok: 0, Message: "Sorry!insufficiet or invalid information"})
						return
					}
					curitem.ID = gotId
					curitem.TotalPrice = curitem.PricePerUnit * curitem.Quantity
					item[index] = curitem
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

func main() {
	ID_generator = 0

	m := mux.NewRouter()
	m.HandleFunc("/", homepage).Methods("GET")
	m.HandleFunc("/showitem", showitem).Methods("GET")
	m.HandleFunc("/additem", additem).Methods("POST")
	m.HandleFunc("/deleteitem/{id}", deleteitem).Methods("DELETE")
	m.HandleFunc("/updateitem/{id}", updateitem).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8800", m))
}
