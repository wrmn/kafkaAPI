package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func logWriter(msg string) {
	log, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer log.Close()

	dt := time.Now()

	_, err := log.Write([]byte(dt.Format("01-02-2006 15:04:05 ") + msg + "\n"))

	if err != nil {
		panic(err)
	}
}

func incJSON(err error, w http.ResponseWriter, transaction Transaction) {
	logWriter(fmt.Sprintf("fail to process \n request for %s\n with error %s", transaction, err.Error()))
	w.WriteHeader(405)
	w.Write([]byte("fail to process with message : " + err.Error()))
	return
}

func incISO(err error, w http.ResponseWriter, iso string) {
	logWriter(fmt.Sprintf("fail to process \n request for %s\n with error %s", iso, err.Error()))
	w.WriteHeader(405)
	w.Write([]byte("fail to process with message : " + err.Error()))
	return
}
