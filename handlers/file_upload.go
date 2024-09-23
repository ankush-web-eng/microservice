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

func handleError(w http.ResponseWriter, errMsg string, code int) {
	http.Error(w, errMsg, code)
}

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
	// Limit file size and set API key from header
	r.ParseMultipartForm(10 << 20)
	apikey := r.Header.Get("API_KEY")
	file, _, err := r.FormFile("file")
	if err != nil {
		handleError(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileChan := make(chan string)
	errChan := make(chan error)

	go func() {
		os.MkdirAll("temp-images", os.ModePerm)
		tempFile, err := os.CreateTemp("temp-images", "upload-*.png")
		if err != nil {
			errChan <- err
			return
		}
		defer tempFile.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			errChan <- err
			return
		}

		_, err = tempFile.Write(fileBytes)
		if err != nil {
			errChan <- err
			return
		}

		fileChan <- tempFile.Name()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch user details concurrently with file processing
	var user models.User
	err = config.DB.Session(&gorm.Session{PrepareStmt: false}).WithContext(ctx).Preload("Cloudinary").First(&user, "api_key = ?", apikey).Error
	if err != nil {
		handleError(w, "User not found", http.StatusNotFound)
		return
	}

	select {
	case tempFileName := <-fileChan:
		// Perform file upload to Cloudinary concurrently
		urlChan := make(chan string)
		go func() {
			url, err := config.UploadFileToCloudinaryAsService(tempFileName, user.Cloudinary.CloudName, user.Cloudinary.APIKey, user.Cloudinary.APISecret)
			if err != nil {
				errChan <- err
				return
			}
			urlChan <- url
		}()

		// Start a GORM transaction for updating requests count
		tx := config.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Select to wait for either Cloudinary upload completion or error
		select {
		case url := <-urlChan:
			// Update the request count for the user in the transaction
			if err := tx.Model(&user.Cloudinary).Update("Requests", gorm.Expr("requests + ?", 1)).Error; err != nil {
				tx.Rollback()
				handleError(w, "Failed to update user", http.StatusInternalServerError)
				return
			}

			// Commit the transaction
			if err := tx.Commit().Error; err != nil {
				handleError(w, "Failed to commit transaction", http.StatusInternalServerError)
				return
			}

			// Clean up the temporary file
			defer os.Remove(tempFileName)

			// Respond with the Cloudinary URL
			w.Write([]byte(url))

		case err := <-errChan:
			tx.Rollback()
			handleError(w, "Failed to upload to Cloudinary: "+err.Error(), http.StatusInternalServerError)

		case <-ctx.Done():
			tx.Rollback()
			handleError(w, "Request timed out", http.StatusRequestTimeout)
		}

	case err := <-errChan:
		if err != nil {
			handleError(w, "Failed to process file: "+err.Error(), http.StatusInternalServerError)
		}

	case <-ctx.Done():
		handleError(w, "Request timed out", http.StatusRequestTimeout)
	}
}
