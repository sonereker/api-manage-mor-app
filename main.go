package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	. "github.com/sonereker/api-manage-mor-app/handler"
	"log"
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

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(":4000", handlers.LoggingHandler(os.Stdout, handlers.CORS(headers, origins, methods)(r))))
}
