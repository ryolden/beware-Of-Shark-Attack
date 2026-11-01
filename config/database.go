package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// 1. Baca dari DB_* atau fallback ke POSTGRES_* bawaan Vercel/Neon
	dbHost := getEnv("DB_HOST", os.Getenv("POSTGRES_HOST"))
	dbUser := getEnv("DB_USER", os.Getenv("POSTGRES_USER"))
	dbPassword := getEnv("DB_PASSWORD", os.Getenv("POSTGRES_PASSWORD"))
	dbName := getEnv("DB_NAME", os.Getenv("POSTGRES_DATABASE"))
	dbPort := getEnv("DB_PORT", os.Getenv("POSTGRES_PORT"))

	// 2. Default jika running lokal tanpa .env
	if dbHost == "" {
		dbHost = "localhost"
		dbUser = "postgres"
		dbPassword = "120126Rop_"
		dbName = "shark_attacks"
		dbPort = "5432"
	}

	if dbPort == "" {
		dbPort = "5432"
	}

	// 3. Tentukan sslmode (Neon Cloud Wajib sslmode=require)
	sslMode := "disable"
	if dbHost != "localhost" && dbHost != "127.0.0.1" {
		sslMode = "require"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Gagal terhubung ke PostgreSQL: %v", err)
	} else {
		fmt.Println("Koneksi PostgreSQL berhasil!")
	}

	return db
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
