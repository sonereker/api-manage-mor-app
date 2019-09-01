package handler

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/sonereker/api-manage-mor-app/model"
	"net/http"
)

func Register(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		employee := model.Account{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&employee); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		if err := db.Save(&employee).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, employee)
		// here, you have access to the db
		// ...
	}
}
