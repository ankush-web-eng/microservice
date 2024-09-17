package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ankush-web-eng/microservice/config"
	email "github.com/ankush-web-eng/microservice/emails"
	"github.com/ankush-web-eng/microservice/models"
	"gorm.io/gorm"
)

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SendEmailRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = email.SendEmail(email.EmailDetails{
		From:    req.From,
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
	})

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}

func SendServiceMailHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SendEmailRequest
	apikey := r.Header.Get("API_KEY")
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if apikey == "" {
		log.Print("APIKEY is required")
		http.Error(w, "API KEY is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var user models.User

	err = config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Preload("Mail").Where("api_key = ?", apikey).First(&user).Error
	if err != nil {
		log.Printf("Failed to find user: %v", err)
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	err = email.SendEmailAsService(email.EmailDetailsAsService{
		From:     req.From,
		To:       req.To,
		Subject:  req.Subject,
		Body:     req.Body,
		Username: user.Mail.Email,
		Password: user.Mail.Password,
	})

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Email sent successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
