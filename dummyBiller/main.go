package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/biller", getBiller).Methods("GET")

	http.ListenAndServe(":5052", r)
}
