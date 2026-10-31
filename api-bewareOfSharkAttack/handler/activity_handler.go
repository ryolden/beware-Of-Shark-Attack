package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"api-bewareOfSharkAttack/entity"
	"api-bewareOfSharkAttack/service"
)

type ActivityHandler struct {
	service *service.ActivityService
}

func NewActivityHandler(service *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

func (h *ActivityHandler) HandleActivities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Menangani ID jika ada pada path: /api/activities/{id}
	var id int
	var err error
	if len(pathParts) > 2 {
		id, err = strconv.Atoi(pathParts[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "ID tidak valid"})
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			act, err := h.service.GetByID(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(act)
		} else {
			activities, err := h.service.GetAll()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(activities)
		}

	case http.MethodPost:
		var act entity.Activity
		if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Payload JSON tidak valid"})
			return
		}
		if err := h.service.Create(&act); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(act)

	case http.MethodPut:
		if id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "ID dibutuhkan untuk update"})
			return
		}
		var act entity.Activity
		if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Payload JSON tidak valid"})
			return
		}
		if err := h.service.Update(id, &act); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(act)

	case http.MethodDelete:
		if id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "ID dibutuhkan untuk hapus"})
			return
		}
		if err := h.service.Delete(id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Berhasil menghapus aktivitas"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method tidak diizinkan"})
	}
}
