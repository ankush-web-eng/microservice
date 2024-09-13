package routes

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ankush-web-eng/microservice/config"
	email "github.com/ankush-web-eng/microservice/emails"
	"github.com/ankush-web-eng/microservice/models"
	"golang.org/x/crypto/bcrypt"
)

func generateOTP() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	user.VerifyCode = strconv.Itoa(generateOTP())

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Error creating user in the database", http.StatusInternalServerError)
		return
	}

	email.SendEmail(email.EmailDetails{
		From:    "ankushsingh.dev@gmail.com",
		To:      []string{user.Email},
		Subject: "Verify your email",
		Body:    "Your verification code is " + user.VerifyCode,
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var input models.User
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "User does not exist", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Logged in"})
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email      string
		VerifyCode string
	}
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.VerifyCode != input.VerifyCode {
		http.Error(w, "Invalid verification code", http.StatusUnauthorized)
		return
	}

	user.IsVerified = true
	config.DB.Save(&user)

	json.NewEncoder(w).Encode(map[string]string{"message": "Account verified"})
}
