package main

import (
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
