package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func postBiller(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logWriter("new Biller request")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logWriter(err.Error())
		fmt.Printf("Error while processing request : %s\n", err.Error())
	}

	var transaction Transaction

	err = json.Unmarshal(b, &transaction)
	if err != nil {
		incJSON(err, w, transaction)
		return
	}

}

func incJSON(err error, w http.ResponseWriter, transaction Transaction) {
	logWriter(fmt.Sprintf("incorrect format \n request for %s\n with error %s", transaction, err.Error()))
	w.WriteHeader(405)
	w.Write([]byte("incorrect format"))
	return
}
