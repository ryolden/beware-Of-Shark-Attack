package handler

import (
	"api-bewareOfSharkAttack/config"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success", "message": "API Beware of Shark Attack (Vercel Serverless)"}`))
	})

	mux.ServeHTTP(w, r)
}
