package routes

import (
	"github.com/ankush-web-eng/microservice/handlers"
	"github.com/gorilla/mux"
)

var CredentialsRoutes = func(router *mux.Router) {

	router.HandleFunc("/apikey", handlers.ApiKeyHandler).Methods("POST")
	router.HandleFunc("/cloudinary", handlers.CloudinaryHanlder).Methods("POST")
	router.HandleFunc("/mail", handlers.MailHandler).Methods("POST")
}
