package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"api-bewareOfSharkAttack/config"
	"api-bewareOfSharkAttack/handler"
	"api-bewareOfSharkAttack/middleware"
	"api-bewareOfSharkAttack/repository"
	"api-bewareOfSharkAttack/service"
	"api-bewareOfSharkAttack/utils"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	// Dependency Injection - User & Auth
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Dependency Injection - Master Activities
	actRepo := repository.NewActivityRepository(db)
	actService := service.NewActivityService(actRepo)
	actHandler := handler.NewActivityHandler(actService)

	// Dependency Injection - Incidents (Tabel Utama)
	incRepo := repository.NewIncidentRepository(db)
	incService := service.NewIncidentService(incRepo)
	incHandler := handler.NewIncidentHandler(incService)

	// Dependency Injection - Seeder CSV Hiu
	seederService := service.NewSeederService(db)

	mux := http.NewServeMux()

	// -------------------- PUBLIC ROUTES --------------------
	mux.HandleFunc("/api/users/register", userHandler.Register)
	mux.HandleFunc("/api/users/login", userHandler.Login)
	mux.HandleFunc("/api/incidents", incHandler.HandleIncidents)
	mux.HandleFunc("/api/incidents/", incHandler.HandleIncidents)

	// -------------------- PROTECTED ROUTES (JWT) --------------------
	mux.HandleFunc("/api/profile", middleware.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Akses area terproteksi JWT berhasil",
			"user":    claims,
		})
	}))

	// CRUD Master Activities (JWT Protected)
	mux.HandleFunc("/api/activities", middleware.JWTMiddleware(actHandler.HandleActivities))
	mux.HandleFunc("/api/activities/", middleware.JWTMiddleware(actHandler.HandleActivities))

	// -------------------- PROTECTED ROUTES (BASIC AUTH) --------------------
	mux.HandleFunc("/api/admin/dashboard", middleware.BasicAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Akses admin terproteksi Basic Auth berhasil"}`))
	}))

	mux.HandleFunc("/api/admin/seed-csv", middleware.BasicAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"message": "Method tidak diizinkan, gunakan POST"})
			return
		}

		filePath := "beware of shark attack(1) - 01 shark attacks lite.csv"
		count, err := seederService.SeedCSV(filePath)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Server berjalan di http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
