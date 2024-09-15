package handlers

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ankush-web-eng/microservice/config"
	"github.com/ankush-web-eng/microservice/models"
	"gorm.io/gorm"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll("temp-images", os.ModePerm)

	tempFile, err := os.CreateTemp("temp-images", "upload-*.png")
	if err != nil {
		http.Error(w, "Unable to create temp file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file bytes", http.StatusInternalServerError)
		return
	}

	tempFile.Write(fileBytes)

	url, err := config.UploadFileToCloudinary(tempFile.Name())
	if err != nil {
		http.Error(w, "Failed to upload to Cloudinary", http.StatusInternalServerError)
		return
	}

	defer os.Remove(tempFile.Name())

	w.Write([]byte(url))
}

func UploadServiceFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	apikey := r.Header.Get("API_KEY")
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll("temp-images", os.ModePerm)

	tempFile, err := os.CreateTemp("temp-images", "upload-*.png")
	if err != nil {
		http.Error(w, "Unable to create temp file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file bytes", http.StatusInternalServerError)
		return
	}

	tempFile.Write(fileBytes)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Preload("Cloudinary").First(&user, "api_key = ?", apikey).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	url, err := config.UploadFileToCloudinaryAsService(tempFile.Name(), user.Cloudinary.CloudName, user.Cloudinary.APIKey, user.Cloudinary.APISecret)
	if err != nil {
		http.Error(w, "Failed to upload to Cloudinary", http.StatusInternalServerError)
		return
	}

	defer os.Remove(tempFile.Name())

	w.Write([]byte(url))
}
