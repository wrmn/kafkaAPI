package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"encoding/json"
)

func sendAPI(data string) string {
	reqBody, err := json.Marshal(data)
	if err != nil {
		logWriter(err.Error())
		panic(err)
	}
	body := bytes.NewReader(reqBody)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://9d32f843-12f1-4402-8d2b-008743bb5e75.mock.pstmn.io/biller", body)
	//req.Header.Set("x-mock-match-request-body", "true")
	res, err := client.Do(req)

	if err != nil {
		logWriter(err.Error())
		panic(err)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logWriter(err.Error())
		panic(err)
	}
	fmt.Printf(string(resBody))

	var resp PaymentResponse

	err = json.Unmarshal(resBody, &resp)

	if err != nil {
		logWriter(err.Error())
		panic(err)
	}
	iso, err := fromJSON(resp)
	if err != nil {
		logWriter(err.Error())
		panic(err)
	}
	return iso
}
