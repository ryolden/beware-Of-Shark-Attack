package middleware

import (
	"net/http"
	"os"
)

func BasicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminUser := os.Getenv("BASIC_AUTH_USER")
		adminPass := os.Getenv("BASIC_AUTH_PASS")
		if adminUser == "" {
			adminUser = "admin"
			adminPass = "admin123"
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != adminUser || pass != adminPass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"status": "error", "message": "Autentikasi Basic Auth gagal"}`))
			return
		}

		next.ServeHTTP(w, r)
	}
}
