package httperror

import (
	"fmt"
	"log"
	"net/http"
)

// Response with given status code
func Error(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

// Response with Bad Request (400)
func BadRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// Response with Internal Server Error (500)
func InternalServerError(w http.ResponseWriter, err error) {
	log.Printf("[InternalServerError] %+v", err)

	Error(w, http.StatusInternalServerError)
}

func LostObject(w http.ResponseWriter, object string) {
	InternalServerError(w, fmt.Errorf("Unexpectedly lost"+object))
}
