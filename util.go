package main

import (
	"net/http"
)

// type ClientMessage string
//
// func RespondWithError(w http.ResponseWriter, msg ClientMessage, status int) {
// 	http.Error(w, string(msg), status)
// }

func RespondWithError(w http.ResponseWriter, msg string, status int) {
	http.Error(w, msg, status)
}
