package routes

import (
	"github.com/ankush-web-eng/microservice/handlers"
	"github.com/gorilla/mux"
)

var AuthRoutes = func(router *mux.Router) {
	router.HandleFunc("/signup", handlers.SignupHandler).Methods("POST")
	router.HandleFunc("/signin", handlers.SigninHandler).Methods("POST")
	router.HandleFunc("/verify", handlers.VerifyHandler).Methods("POST")
	router.HandleFunc("/signin/verify", handlers.AuthVerifier).Methods("POST")
	router.HandleFunc("/upload", handlers.UploadFileHandler).Methods("POST")
	router.HandleFunc("/send-email", handlers.SendEmailHandler).Methods("POST")
}
