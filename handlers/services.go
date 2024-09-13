package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ankush-web-eng/microservice/config"
	"github.com/ankush-web-eng/microservice/helpers"
	"github.com/ankush-web-eng/microservice/models"
	"gorm.io/gorm"
)

func ApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	apikey, err := helpers.GenerateAPIKey()
	if err != nil {
		http.Error(w, "Error in generating API Key", http.StatusInternalServerError)
		return
	}

	var user models.User

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	txErr := config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("email = ?", input.Email).First(&user).Error
		if err != nil {
			log.Printf("Error in fetching user: %v", err)
			return err
		}
		user.APIKey = &apikey
		return tx.Save(&user).Error
	})

	if txErr != nil {
		log.Printf("Error in updating user with API Key: %v", txErr)
		http.Error(w, "Error in updating user with API Key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"apikey": apikey})
}
