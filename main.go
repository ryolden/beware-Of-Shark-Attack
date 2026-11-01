package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"

	"api-bewareOfSharkAttack/config"
	"api-bewareOfSharkAttack/handler"
	"api-bewareOfSharkAttack/middleware"
	"api-bewareOfSharkAttack/repository"
	"api-bewareOfSharkAttack/service"
	"api-bewareOfSharkAttack/utils"
)

//go:embed api/data.csv
var csvData []byte

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment OS/Vercel")
	}

	db := config.ConnectDB()
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	actRepo := repository.NewActivityRepository(db)
	actService := service.NewActivityService(actRepo)
	actHandler := handler.NewActivityHandler(actService)

	incRepo := repository.NewIncidentRepository(db)
	incService := service.NewIncidentService(incRepo)
	incHandler := handler.NewIncidentHandler(incService)

	seederService := service.NewSeederService(db)

	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/api/users/register", userHandler.Register)
	mux.HandleFunc("/api/users/login", userHandler.Login)
	mux.HandleFunc("/api/incidents", incHandler.HandleIncidents)
	mux.HandleFunc("/api/incidents/", incHandler.HandleIncidents)

	// Protected Routes (JWT)
	mux.HandleFunc("/api/profile", middleware.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Akses area terproteksi JWT berhasil",
			"user":    claims,
		})
	}))

	mux.HandleFunc("/api/activities", middleware.JWTMiddleware(actHandler.HandleActivities))
	mux.HandleFunc("/api/activities/", middleware.JWTMiddleware(actHandler.HandleActivities))

	// Protected Routes (Basic Auth)
	mux.HandleFunc("/api/admin/dashboard", middleware.BasicAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Akses admin terproteksi Basic Auth berhasil"}`))
	}))

	mux.HandleFunc("/api/admin/seed-csv", middleware.BasicAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		reader := bytes.NewReader(csvData)
		count, err := seederService.SeedCSV(reader)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Import CSV data insiden berhasil",
			"total":   count,
		})
	}))

	// Root Route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "error",
				"message": "Endpoint tidak ditemukan",
				"path":    r.URL.Path,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success", "message": "API Beware of Shark Attack is Live"}`))
	})

	// Middleware/Handler untuk memulihkan Path dari X-Forwarded-Uri Vercel
	handlerWithRewrite := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if forwardedURI := r.Header.Get("X-Forwarded-Uri"); forwardedURI != "" {
			if parsedURL, err := url.Parse(forwardedURI); err == nil {
				r.URL.Path = parsedURL.Path
			}
		}
		mux.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Server berjalan di port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithRewrite))
}
