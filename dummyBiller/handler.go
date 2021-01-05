package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getBiller(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	var transactions transaction

	//s, _ := strconv.Unquote(string(b))

	err = json.Unmarshal([]byte(b), &transactions)
	if err != nil {
		fmt.Println(err.Error())
	}

	var response paymentResponse

	response.TransactionData = transactions
	response.ResponseStatus.ResponseCode = 200
	response.ResponseStatus.ReasonCode = 0
	response.ResponseStatus.ResponseDescription = "success"

	resJson, err := json.MarshalIndent(response, "", "   ")
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(resJson)
}
