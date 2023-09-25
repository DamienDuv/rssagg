package api

import "net/http"

func HandlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Somthing went wrong")
}