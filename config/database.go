package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// 1. Cek DATABASE_URL / POSTGRES_URL dari Neon/Vercel terlebih dahulu
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("POSTGRES_URL")
	}

	var db *sql.DB
	var err error

	if dbURL != "" {
		db, err = sql.Open("postgres", dbURL)
	} else {
		// Fallback ke variabel individu / lokal
		dbHost := getEnv("DB_HOST", os.Getenv("POSTGRES_HOST"))
		dbUser := getEnv("DB_USER", os.Getenv("POSTGRES_USER"))
		dbPassword := getEnv("DB_PASSWORD", os.Getenv("POSTGRES_PASSWORD"))
		dbName := getEnv("DB_NAME", os.Getenv("POSTGRES_DATABASE"))
		dbPort := getEnv("DB_PORT", "5432")

		if dbHost == "" {
			dbHost = "localhost"
			dbUser = "postgres"
			dbPassword = "120126Rop_"
			dbName = "shark_attacks"
		}

		sslMode := "disable"
		if dbHost != "localhost" && dbHost != "127.0.0.1" {
			sslMode = "require"
		}

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
		db, err = sql.Open("postgres", dsn)
	}

	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	if err = db.Ping(); err != nil {
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
