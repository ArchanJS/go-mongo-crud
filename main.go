package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()
	s := route.PathPrefix("/api").Subrouter()

	// routes
	s.HandleFunc("/createuser", createUser).Methods("POST")
	s.HandleFunc("/getuser", getUsers).Methods("GET")
	s.HandleFunc("/updateuser", updateUser).Methods("PUT")
	s.HandleFunc("/deleteuser/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", s))
}
