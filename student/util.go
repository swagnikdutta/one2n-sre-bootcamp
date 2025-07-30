package student

import (
	"net/http"
)

func RespondWithError(w http.ResponseWriter, msg string, status int) {
	http.Error(w, msg, status)
}
