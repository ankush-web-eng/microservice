package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ankush-web-eng/microservice/handlers"
	"github.com/gorilla/mux"
)

var ServiceRoutes = func(router *mux.Router) {
	router.HandleFunc("/", testRoute).Methods("GET")
	router.HandleFunc("/send-mail", handlers.SendServiceMailHandler).Methods("POST")
	router.HandleFunc("/upload-file", handlers.UploadServiceFileHandler).Methods("POST")
}

func testRoute(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "This is a simple service developed by Ankush"})
}
