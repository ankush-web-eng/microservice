package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/ankush-web-eng/microservice/config"
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

	w.Write([]byte(url))
}
