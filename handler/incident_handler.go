package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"api-bewareOfSharkAttack/entity"
	"api-bewareOfSharkAttack/middleware"
	"api-bewareOfSharkAttack/service"
	"api-bewareOfSharkAttack/utils"
)

type IncidentHandler struct {
	service *service.IncidentService
}

func NewIncidentHandler(service *service.IncidentService) *IncidentHandler {
	return &IncidentHandler{service: service}
}

func (h *IncidentHandler) HandleIncidents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	var id int
	if len(pathParts) > 2 {
		id, _ = strconv.Atoi(pathParts[2])
	}

	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			inc, err := h.service.GetByID(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(inc)
		} else {
			limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
			page, _ := strconv.Atoi(r.URL.Query().Get("page"))

			incidents, err := h.service.GetAll(limit, page)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(incidents)
		}

	case http.MethodPost:
		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		var inc entity.Incident
		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Payload JSON tidak valid"})
			return
		}

		if ok && claims != nil {
			inc.UserID = &claims.UserID
		}

		if err := h.service.Create(&inc); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(inc)

	case http.MethodDelete:
		if id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "ID insiden dibutuhkan"})
			return
		}

		if err := h.service.Delete(id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Berhasil menghapus data insiden"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method tidak diizinkan"})
	}
}
