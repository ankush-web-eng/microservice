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

func handleJSONError(w http.ResponseWriter, errMsg string, code int) {
	w.WriteHeader(code)
	response := map[string]string{"error": errMsg}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

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
		handleJSONError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if apikey == "" {
		log.Print("APIKEY is required")
		handleJSONError(w, "API KEY is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Find user based on API key
	var user models.User
	err = config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Preload("Mail").Where("api_key = ?", apikey).First(&user).Error
	if err != nil {
		log.Printf("Failed to find user: %v", err)
		handleJSONError(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	// Goroutine for sending email asynchronously
	emailChan := make(chan error)
	go func() {
		err = email.SendEmailAsService(email.EmailDetailsAsService{
			From:     req.From,
			To:       req.To,
			Subject:  req.Subject,
			Body:     req.Body,
			Username: user.Mail.Email,
			Password: user.Mail.Password,
		})
		emailChan <- err
	}()

	// Goroutine for updating mail requests asynchronously
	updateChan := make(chan error)
	go func() {
		tx := config.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		err = tx.Model(&models.Mail{}).Where("user_id = ?", user.ID).Updates(models.Mail{Requests: user.Mail.Requests + 1}).Error
		if err != nil {
			tx.Rollback()
			updateChan <- err
			return
		}

		if err := tx.Commit().Error; err != nil {
			updateChan <- err
			return
		}
		updateChan <- nil
	}()

	// Wait for the email and DB update operations to complete
	select {
	case emailErr := <-emailChan:
		if emailErr != nil {
			log.Printf("Failed to send email: %v", emailErr)
			handleJSONError(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

	case updateErr := <-updateChan:
		if updateErr != nil {
			log.Printf("Failed to update mail requests: %v", updateErr)
			handleJSONError(w, "Failed to update mail requests", http.StatusInternalServerError)
			return
		}

	case <-ctx.Done():
		handleJSONError(w, "Request timed out", http.StatusRequestTimeout)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Email sent successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		handleJSONError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
