package main

import (
	"fmt"
	"net/http"

	"github.com/ankush-web-eng/microservice/config"
	"github.com/ankush-web-eng/microservice/controllers"
	"github.com/ankush-web-eng/microservice/handlers"
	"github.com/ankush-web-eng/microservice/middlewares"
	"github.com/ankush-web-eng/microservice/models"
	"github.com/ankush-web-eng/microservice/routes"
	"github.com/gorilla/mux"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	config.InitDB()
	config.InitCloudinary()
	config.DB.AutoMigrate(&models.User{})

	router := mux.NewRouter()

	router.HandleFunc("/signup", routes.SignupHandler).Methods("POST")
	router.HandleFunc("/signin", routes.SigninHandler).Methods("POST")
	router.HandleFunc("/verify", routes.VerifyHandler).Methods("POST")
	router.HandleFunc("/upload", handlers.UploadFileHandler).Methods("POST")

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Use(middlewares.JWTAuthMiddleware)
	authRouter.HandleFunc("/protected", controllers.ProtectedHandler).Methods("GET")

	http.ListenAndServe(":8080", router)
}
