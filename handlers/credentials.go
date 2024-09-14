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
	json.NewEncoder(w).Encode(apikey)
}

func CloudinaryHanlder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string `json:"email"`
		CloudName string `json:"cloudname"`
		APIKey    string `json:"apikey"`
		APISecret string `json:"apisecret"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user models.User
		err := tx.Where("email = ?", input.Email).First(&user).Error
		if err != nil {
			log.Printf("Error in fetching user: %v", err)
			return err
		}

		createErr := tx.Create(&models.Cloudinary{
			APIKey:    input.APIKey,
			APISecret: input.APISecret,
			CloudName: input.CloudName,
			UserID:    user.ID,
		}).Error
		if createErr != nil {
			log.Printf("Error in creating cloudinary: %v", createErr)
			return createErr
		}
		return nil
	})

	if err != nil {
		log.Printf("Error in creating cloudinary: %v", err)
		http.Error(w, "Error in creating cloudinary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cloudinary configuration saved successfully"})
}

func MailHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EmailUser string `json:"emailuser"`
		Password  string `json:"password"`
		Email     string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	dbErr := config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user models.User
		err := tx.Where("email = ?", input.Email).First(&user).Error
		if err != nil {
			log.Printf("Error in fetching user: %v", err)
			return err
		}

		createErr := tx.Create(&models.Mail{
			Email:    input.EmailUser,
			Password: input.Password,
			UserID:   user.ID,
		}).Error
		if createErr != nil {
			log.Printf("Error in creating mail: %v", createErr)
			return createErr
		}
		return nil
	})

	if dbErr != nil {
		log.Printf("Error in creating mail: %v", dbErr)
		http.Error(w, "Error in creating mail", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Mail configuration saved successfully"})
}
