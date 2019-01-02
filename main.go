package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	. "github.com/sonereker/api-manage-mor-app/handler"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	r.Handle("/authenticate", AuthenticationHandler).Methods("POST")
	r.Handle("/collection", ValidateToken(ViewCollectionHandler)).Methods("GET")
	r.Handle("/collection", ValidateToken(AddWorkHandler)).Methods("POST")
	r.Handle("/collection/work/{id}", ValidateToken(UpdateWorkHandler)).Methods("PUT")
	r.Handle("/collection/work/{id}", ValidateToken(DeleteWorkHandler)).Methods("DELETE")

	_ = http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}
