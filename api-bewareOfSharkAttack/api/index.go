package handler

import (
	"encoding/json"
	"net/http"

	"api-bewareOfSharkAttack/config"
	"api-bewareOfSharkAttack/handler"
	"api-bewareOfSharkAttack/middleware"
	"api-bewareOfSharkAttack/repository"
	"api-bewareOfSharkAttack/service"
	"api-bewareOfSharkAttack/utils"
)

// PERHATIKAN: Nama fungsi wajib Handler (H kapital)
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

	// Root route agar url utama tidak 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success", "message": "API Beware of Shark Attack is Live on Vercel"}`))
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
		filePath := "beware of shark attack(1) - 01 shark attacks lite.csv"
		count, err := seederService.SeedCSV(filePath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"total":  count,
		})
	}))

	// Teruskan request HTTP ke ServeMux
	mux.ServeHTTP(w, r)
}
