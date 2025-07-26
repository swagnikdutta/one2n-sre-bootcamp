package main

import (
	"fmt"
	"log"
	"net/http"
)

type LogMessage string
type ClientMessage string

func RespondWithError(w http.ResponseWriter, logMsg LogMessage, clientMsg ClientMessage, status int) {
	log.Println(logMsg)
	fmt.Println(logMsg)
	http.Error(w, string(clientMsg), status)
}
