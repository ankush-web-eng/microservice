package config

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var Cloudinary *cloudinary.Cloudinary

func InitCloudinary() {
	var err error
	Cloudinary, err = cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}
}

func UploadFileToCloudinary(filePath string) (string, error) {
	uploadResult, err := Cloudinary.Upload.Upload(context.Background(), filePath, uploader.UploadParams{})
	if err != nil {
		log.Fatal("Failed to upload file to Cloudinary:", err)
		return "", err
	}
	return uploadResult.SecureURL, nil
}
