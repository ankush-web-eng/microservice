package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ankush-web-eng/microservice/config"
	"github.com/ankush-web-eng/microservice/handlers"
	"github.com/ankush-web-eng/microservice/models"
	"github.com/ankush-web-eng/microservice/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	config.InitDB()
	config.DB.AutoMigrate(&models.User{})
	config.DB.AutoMigrate(&models.Mail{})
	config.DB.AutoMigrate(&models.Cloudinary{})

	router := mux.NewRouter()

	credentialsRouter := router.PathPrefix("/credentials").Subrouter()
	serviceRouter := router.PathPrefix("/service").Subrouter()

	routes.AuthRoutes(router)
	routes.CredentialsRoutes(credentialsRouter)
	routes.ServiceRoutes(serviceRouter)

	router.HandleFunc("/upload", handlers.UploadFileHandler).Methods("POST")
	router.HandleFunc("/send-email", handlers.SendEmailHandler).Methods("POST")

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "API_KEY"},
		AllowCredentials: true,
	})

	handler := corsOptions.Handler(router)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
