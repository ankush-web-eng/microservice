package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ankush-web-eng/microservice/config"
	email "github.com/ankush-web-eng/microservice/emails"
	"github.com/ankush-web-eng/microservice/helpers"
	"github.com/ankush-web-eng/microservice/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func sendVerificationEmail(userEmail, verifyCode string) {
	go func() {
		err := email.SendEmail(email.EmailDetails{
			From:    "ankushsingh.dev@gmail.com",
			To:      []string{userEmail},
			Subject: "Verify your email",
			Body:    "Your verification code is " + verifyCode,
		})
		if err != nil {
			log.Printf("Error sending email: %v", err)
		}
	}()
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingUser models.User
		if err := tx.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			if existingUser.IsVerified {
				return fmt.Errorf("user already exists")
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("error while hashing password: %w", err)
			}
			existingUser.Password = string(hashedPassword)

			otp, err := helpers.GenerateOTP()
			if err != nil {
				return fmt.Errorf("error generating OTP: %w", err)
			}
			existingUser.VerifyCode = strconv.Itoa(otp)

			if err := tx.Save(&existingUser).Error; err != nil {
				return fmt.Errorf("error updating user in the database: %w", err)
			}

			user = existingUser
		} else if err != gorm.ErrRecordNotFound {
			return err
		} else {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

			if err != nil {
				return fmt.Errorf("error while hashing password: %w", err)
			}
			user.Password = string(hashedPassword)

			otp, err := helpers.GenerateOTP()
			if err != nil {
				return fmt.Errorf("error generating OTP: %w", err)
			}
			user.VerifyCode = strconv.Itoa(otp)

			if err := tx.Create(&user).Error; err != nil {
				return fmt.Errorf("error creating user in the database: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error in SignupHandler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendVerificationEmail(user.Email, user.VerifyCode)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Printf("Error in SigninHandler: %v", err)
		http.Error(w, "User does not exist", http.StatusUnauthorized)
		return
	}

	if !user.IsVerified {
		log.Printf("User %s is not verified", user.Email)
		http.Error(w, "User is not verified", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("Password mismatch: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("User %s signed in successfully", user.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email      string
		VerifyCode string
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Printf("Error in VerifyHandler: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.VerifyCode != input.VerifyCode {
		http.Error(w, "Invalid verification code", http.StatusUnauthorized)
		return
	}

	user.IsVerified = true
	if err := config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Save(&user).Error; err != nil {
		log.Printf("Error in VerifyHandler: %v", err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Account verified"})
}

func AuthVerifier(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Printf("Error in AuthVerifier: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"isVerified": user.IsVerified})
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	var user models.User

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	done := make(chan error)
	go func() {
		err := config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Preload("Cloudinary").Preload("Mail").Where("email = ?", email).First(&user).Error
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Error in GetUserHandler: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
