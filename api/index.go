package handler

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"

	"api-bewareOfSharkAttack/config"
	"api-bewareOfSharkAttack/handler"
	"api-bewareOfSharkAttack/middleware"
	"api-bewareOfSharkAttack/repository"
	"api-bewareOfSharkAttack/service"
	"api-bewareOfSharkAttack/utils"
)

//go:embed data.csv
var csvData []byte

func Handler(w http.ResponseWriter, r *http.Request) {
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

	// Handler utama (Root / Default)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Izinkan request ke /, empty, atau /api/index.go
		if r.URL.Path == "/" || r.URL.Path == "" || r.URL.Path == "/api/index.go" {
			// Ambil 10 data insiden pertama dari Neon PostgreSQL
			incidents, err := incService.GetAll(10, 1)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "success",
				"total":  len(incidents),
				"data":   incidents,
			})
			return
		}

		// Path tidak dikenal
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Endpoint tidak ditemukan",
			"path":    r.URL.Path,
		})
	})

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

	// Protected Routes (Basic Auth & Seeder)
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

	mux.ServeHTTP(w, r)
}
