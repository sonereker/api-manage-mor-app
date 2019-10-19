package handler

import (
	"encoding/json"
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
