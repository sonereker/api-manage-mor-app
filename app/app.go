package app

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sonereker/api-manage-mor-app/config"
	"github.com/sonereker/api-manage-mor-app/handler"
	"github.com/sonereker/api-manage-mor-app/model"
	"log"
	"net/http"
	"os"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(config *config.Config) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.SSLMode,
	)

	db, err := gorm.Open(config.DB.Dialect, dbInfo)
	if err != nil {
		log.Fatal("Could not connect database: ", err)
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	v1 := a.Router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/authenticate", handler.Authenticate).Methods("POST")
	v1.HandleFunc("/accounts", handler.Register(a.DB)).Methods("POST")
	c := v1.PathPrefix("/collection").Subrouter()
	c.HandleFunc("/", handler.ValidateToken(handler.ViewCollection)).Methods("GET")
}

func (a *App) Run(host string) {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(host, handlers.LoggingHandler(os.Stdout, handlers.CORS(headers, origins, methods)(a.Router))))
}
