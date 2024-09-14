package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ankush-web-eng/microservice/handlers"
	"github.com/gorilla/mux"
)

var ServiceHandler = func(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "This is a sample microservice"})
	}).Methods("GET")
	router.HandleFunc("/service/apikey", handlers.ApiKeyHandler).Methods("POST")
	router.HandleFunc("/service/cloudinary", handlers.CloudinaryHanlder).Methods("POST")
	router.HandleFunc("/service/mail", handlers.MailHandler).Methods("POST")
}
