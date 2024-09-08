package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/ankush-web-eng/microservice/utils"
	"github.com/dgrijalva/jwt-go"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.Split(authHeader, "Bearer ")[1]
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return
		}

		type contextKey string

		ctx := context.WithValue(r.Context(), contextKey("userID"), claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
