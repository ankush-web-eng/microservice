package controllers

import (
	"encoding/json"
	"net/http"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	response := map[string]interface{}{
		"message": "Welcome to the protected route!",
		"userID":  userID,
	}
	json.NewEncoder(w).Encode(response)
}
