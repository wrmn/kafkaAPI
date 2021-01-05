package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"encoding/json"
)

func sendAPI(data string) string {
	reqBody, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		logWriter(err.Error())
		fmt.Println(err.Error())
		return ""
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://tiruan.herokuapp.com/biller", bytes.NewBuffer(reqBody))
	//req.Header.Set("x-mock-match-request-body", "true")
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logWriter(err.Error())
		fmt.Println(err.Error())
		return ""
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logWriter(err.Error())

		fmt.Println(err.Error())
		return ""
	}
	fmt.Print(string(resBody) + "wow")

	var resp PaymentResponse

	err = json.Unmarshal(resBody, &resp)

	fmt.Printf("%v", resp)

	if err != nil {
		logWriter(err.Error())
		fmt.Println(err.Error())
		return ""
	}
	iso, err := fromJSON(resp)
	if err != nil {
		logWriter(err.Error())
		fmt.Println(err.Error())
		return ""
	}
	return iso
}
