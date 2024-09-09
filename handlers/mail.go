package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	email "github.com/ankush-web-eng/microservice/emails"
	"github.com/ankush-web-eng/microservice/models"
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
