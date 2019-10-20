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
	dbInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Name,
		config.DB.SSLMode,
	)

	db, err := gorm.Open(config.DB.Dialect, dbInfo)
	if err != nil {
		log.Fatal("Could not connect database: ", err)
	}

	a.DB = DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	v1 := a.Router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/authenticate", handler.Authenticate).Methods("POST")
	v1.HandleFunc("/accounts", handler.Register(a.DB)).Methods("POST")

	c := v1.PathPrefix("/catalogs").Subrouter()
	c.HandleFunc("", handler.ValidateToken(a.GetAllCatalogs)).Methods("GET")
	c.HandleFunc("", handler.ValidateToken(a.CreateCatalog)).Methods("POST")
	c.HandleFunc("/{id}", handler.ValidateToken(a.GetCatalog)).Methods("GET")
	c.HandleFunc("/{id}", handler.ValidateToken(a.UpdateCatalog)).Methods("PUT")
	c.HandleFunc("/{id}", handler.ValidateToken(a.DeleteCatalog)).Methods("DELETE")

	c.HandleFunc("/{id}/assets", handler.ValidateToken(a.GetAllAssets)).Methods("GET")
	c.HandleFunc("/{id}/assets/new", handler.ValidateToken(a.CreateAsset)).Methods("POST")
	c.HandleFunc("/{id}/assets/{assetId}", handler.ValidateToken(a.GetAsset)).Methods("GET")
	c.HandleFunc("/{id}/assets/{assetId}", handler.ValidateToken(a.UpdateAsset)).Methods("PUT")
	c.HandleFunc("/{id}/assets/{assetId}", handler.ValidateToken(a.DeleteAsset)).Methods("DELETE")
}

func (a *App) Run(host string) {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(host, handlers.LoggingHandler(os.Stdout, handlers.CORS(headers, origins, methods)(a.Router))))
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&model.Account{})
	db.AutoMigrate(&model.Catalog{})
	db.AutoMigrate(&model.Asset{})
	return db
}

// Handlers to manage Catalog Data
func (a *App) GetAllCatalogs(w http.ResponseWriter, r *http.Request) {
	handler.GetAllCatalogs(a.DB, w, r)
}

func (a *App) CreateCatalog(w http.ResponseWriter, r *http.Request) {
	handler.CreateCatalog(a.DB, w, r)
}

func (a *App) GetCatalog(w http.ResponseWriter, r *http.Request) {
	handler.GetCatalog(a.DB, w, r)
}

func (a *App) UpdateCatalog(w http.ResponseWriter, r *http.Request) {
	handler.UpdateCatalog(a.DB, w, r)
}

func (a *App) DeleteCatalog(w http.ResponseWriter, r *http.Request) {
	handler.DeleteCatalog(a.DB, w, r)
}

// Handlers to manage Asset Data
func (a *App) GetAllAssets(w http.ResponseWriter, r *http.Request) {
	handler.GetAllAssets(a.DB, w, r)
}

func (a *App) CreateAsset(w http.ResponseWriter, r *http.Request) {
	handler.CreateAsset(a.DB, w, r)
}

func (a *App) GetAsset(w http.ResponseWriter, r *http.Request) {
	handler.GetAsset(a.DB, w, r)
}

func (a *App) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	handler.UpdateAsset(a.DB, w, r)
}

func (a *App) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	handler.DeleteAsset(a.DB, w, r)
}
