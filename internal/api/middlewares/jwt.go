package middlewares

import (
	"fmt"
	"ftm-explorer/internal/logger"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JwtMiddleware defines HTTP handler middleware for JWT authentication.
func JwtMiddleware(next http.Handler, log logger.ILogger, secret string, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			log.Warning("missing authorization token")
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// parse the authorization header
		token, err := jwt.Parse(authHeader[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method '%v'", t.Header["alg"])
			}
			// Return the key for validation
			return []byte(secret), nil
		})
		if err != nil {
			log.Warningf("failed to parse authorization token: %v", err)
			http.Error(w, fmt.Sprintf("Failed to parse authorization token: %v", err), http.StatusUnauthorized)
			return
		}

		// check token is valid
		if !token.Valid {
			log.Warning("invalid authorization token")
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Warning("invalid authorization token")
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// check version is valid
		if claims["version"] != version {
			log.Warning("invalid authorization token version")
			http.Error(w, "Invalid authorization token version", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
