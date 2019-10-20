package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sonereker/api-manage-mor-app/model"
	"net/http"
)

func GetAllCatalogs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var catalogs []model.Catalog
	db.Find(&catalogs)
	respondJSON(w, http.StatusOK, catalogs)
}

func CreateCatalog(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	catalog := model.Catalog{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&catalog); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&catalog).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, catalog)
}

func GetCatalog(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	catalog := getCatalogOr404(db, id, w, r)
	if catalog == nil {
		return
	}
	respondJSON(w, http.StatusOK, catalog)
}

func UpdateCatalog(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	catalog := getCatalogOr404(db, uuid, w, r)
	if catalog == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&catalog); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&catalog).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, catalog)
}

func DeleteCatalog(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	catalog := getCatalogOr404(db, uuid, w, r)
	if catalog == nil {
		return
	}
	if err := db.Delete(&catalog).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func getCatalogOr404(db *gorm.DB, uuid string, w http.ResponseWriter, r *http.Request) *model.Catalog {
	catalog := model.Catalog{}
	if err := db.First(&catalog, uuid).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &catalog
}
