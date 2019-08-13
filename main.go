package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sonereker/api-manage-mor-app/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/authenticate", handlers.AuthenticationHandler).Methods("POST")
	r.Handle("/collection", handlers.ValidateToken(handlers.ViewCollectionHandler)).Methods("GET")
	r.Handle("/collection", handlers.ValidateToken(handlers.AddWorkHandler)).Methods("POST")
	r.Handle("/collection/work/{id}", handlers.ValidateToken(handlers.UpdateWorkHandler)).Methods("PUT")
	r.Handle("/collection/work/{id}", handlers.ValidateToken(handlers.DeleteWorkHandler)).Methods("DELETE")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(":4000", handlers.LoggingHandler(os.Stdout, handlers.CORS(headers, origins, methods)(r))))
}
