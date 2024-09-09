package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ankush-web-eng/microservice/config"
	"github.com/ankush-web-eng/microservice/controllers"
	"github.com/ankush-web-eng/microservice/handlers"
	"github.com/ankush-web-eng/microservice/middlewares"
	"github.com/ankush-web-eng/microservice/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file in railway")
		}
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	config.InitDB()
	config.InitCloudinary()
	// config.DB.AutoMigrate(&models.User{})

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "This is a sample microservice"})
	}).Methods("GET")

	router.HandleFunc("/signup", routes.SignupHandler).Methods("POST")
	router.HandleFunc("/signin", routes.SigninHandler).Methods("POST")
	router.HandleFunc("/verify", routes.VerifyHandler).Methods("POST")
	router.HandleFunc("/upload", handlers.UploadFileHandler).Methods("POST")
	router.HandleFunc("/send-email", handlers.SendEmailHandler).Methods("POST")

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Use(middlewares.JWTAuthMiddleware)
	authRouter.HandleFunc("/protected", controllers.ProtectedHandler).Methods("GET")

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://go.ankushsingh.tech"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
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
