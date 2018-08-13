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
	ID int `json:"id"`
	Name string `json:"name"`
	Quantity int `json:"quantity"`
	Price int `json:"price"`
}
var item []Item

func homepage(w http.ResponseWriter , r *http.Request ) {
   fmt.Fprintf(w, "Welcome to Shopping List RESTful api server")
}

func showitem(w http.ResponseWriter , r *http.Request ) {
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(item)
}

func additem(w http.ResponseWriter , r *http.Request ) {
    w.Header().Set("Content-Type", "application/json")
    var curitem Item
    _ = json.NewDecoder(r.Body).Decode(&curitem)
    ID_generator++
    curitem.ID = ID_generator
    item = append(item , curitem)
    json.NewEncoder(w).Encode(curitem)
}
func deleteitem(w http.ResponseWriter , r *http.Request ) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    //parameters have been handled
    gotId , err := strconv.Atoi(params["id"])
    if err == nil {
	    for index,i := range item {
	    	if i.ID == gotId {
	            item = append(item[:index] , item[index+1 :]...)
	            //found the desired index , removed and then appended the rest of the items items
	            break;
	    	}
	    }
   }
    json.NewEncoder(w).Encode(item)

}
func updateitem(w http.ResponseWriter , r *http.Request ) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    gotId , err := strconv.Atoi(params["id"])
    if err == nil {
	    for index,i := range item {
	    	if i.ID == gotId {
	            item = append(item[:index] , item[index+1 :]...)
	            //basically deletion
			           
			    var curitem Item
			    curitem = i
			    _ = json.NewDecoder(r.Body).Decode(&curitem)
			    item = append(item , curitem)
			    //updated item has been added
	            break;
	    	}
	    }
   }
   json.NewEncoder(w).Encode(item)

}

func main() {
	ID_generator = 1
    //sample items added
     item = append(item , Item{ID:1 , Name:"Fish", Quantity:5 , Price : 350})
     item = append(item , Item{ID:2 , Name:"Rice", Quantity:2 , Price : 120})
    // item = append(item , Item{ID:3 , Name:"fuck", Quantity:5 , Price : 350})
	m := mux.NewRouter()
    //router
    //function handler
	m.HandleFunc("/",  homepage).Methods("GET")
	m.HandleFunc("/showitem", showitem).Methods("GET")	
	m.HandleFunc("/additem", additem).Methods("POST")
    m.HandleFunc("/deleteitem/{id}", deleteitem).Methods("DELETE")
    m.HandleFunc("/updateitem/{id}", updateitem).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8800", m))
}
