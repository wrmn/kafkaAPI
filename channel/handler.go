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

	result, err := fromJSON(transaction)
	if err != nil {
		incISO(err, w, result)
		return
	}

	isoRes := []byte(result)
	w.WriteHeader(200)
	w.Write(isoRes)

}
