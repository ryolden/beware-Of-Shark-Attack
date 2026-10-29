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

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/api/users/register", userHandler.Register)
	mux.HandleFunc("/api/users/login", userHandler.Login)

	// Protected Routes Contoh JWT
	mux.HandleFunc("/api/profile", middleware.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Akses area terproteksi JWT berhasil",
			"user":    claims,
		})
	}))

	// Protected Routes Contoh Basic Auth
	mux.HandleFunc("/api/admin/dashboard", middleware.BasicAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Akses admin terproteksi Basic Auth berhasil"}`))
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
