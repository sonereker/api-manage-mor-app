package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sonereker/api-manage-mor-app/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/authenticate", handler.AuthenticationHandler).Methods("POST")
	r.HandleFunc("/collection", handler.ValidateToken(handler.ViewCollectionHandler)).Methods("GET")
	r.HandleFunc("/collection", handler.ValidateToken(handler.AddWorkHandler)).Methods("POST")
	r.HandleFunc("/collection/work/{id}", handler.ValidateToken(handler.UpdateWorkHandler)).Methods("PUT")
	r.HandleFunc("/collection/work/{id}", handler.ValidateToken(handler.DeleteWorkHandler)).Methods("DELETE")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(":4000", handlers.LoggingHandler(os.Stdout, handlers.CORS(headers, origins, methods)(r))))
}
