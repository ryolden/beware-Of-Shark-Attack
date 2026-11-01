package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// 1. Cek apakah ada DATABASE_URL / POSTGRES_URL dari Neon/Vercel
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = os.Getenv("POSTGRES_URL")
	}

	// 2. Jika tidak ada DATABASE_URL, buat string koneksi dari variabel DB_*
	if connStr == "" {
		dbHost := getEnv("DB_HOST", os.Getenv("POSTGRES_HOST"))
		dbUser := getEnv("DB_USER", os.Getenv("POSTGRES_USER"))
		dbPassword := getEnv("DB_PASSWORD", os.Getenv("POSTGRES_PASSWORD"))
		dbName := getEnv("DB_NAME", os.Getenv("POSTGRES_DATABASE"))
		dbPort := getEnv("DB_PORT", os.Getenv("POSTGRES_PORT"))

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

		sslMode := "disable"
		if dbHost != "localhost" && dbHost != "127.0.0.1" {
			sslMode = "require"
		}

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	}

	// 3. Buka koneksi ke Database
	db, err := sql.Open("postgres", connStr)
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
