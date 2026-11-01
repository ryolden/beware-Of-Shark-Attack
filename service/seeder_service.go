package service

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type SeederService struct {
	DB *sql.DB
}

func NewSeederService(db *sql.DB) *SeederService {
	return &SeederService{DB: db}
}

func (s *SeederService) SeedCSV(r io.Reader) (int, error) {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true

	// Baca header baris pertama
	_, err := reader.Read()
	if err != nil {
		return 0, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Cache In-Memory
	locationMap := make(map[string]int)
	speciesMap := make(map[string]int)
	activityMap := make(map[string]int)

	stmtIncident, err := tx.Prepare(`
		INSERT INTO incidents (incident_date, year, location_id, species_id, activity_id, injury, fatal)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return 0, err
	}
	defer stmtIncident.Close()

	importedCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		dateStr := strings.TrimSpace(record[0])
		yearStr := strings.TrimSpace(record[2])
		countryStr := strings.TrimSpace(record[4])
		areaStr := strings.TrimSpace(record[5])
		locationStr := strings.TrimSpace(record[6])
		activityStr := strings.TrimSpace(record[7])
		injuryStr := strings.TrimSpace(record[13])
		fatalStr := strings.TrimSpace(record[14])
		speciesStr := strings.TrimSpace(record[15])

		if countryStr == "" {
			countryStr = "UNKNOWN"
		}
		if speciesStr == "" {
			speciesStr = "Unidentified Species"
		}
		if activityStr == "" {
			activityStr = "Unknown Activity"
		}

		year, _ := strconv.Atoi(yearStr)

		// --- 1. Get/Insert Location (Aman dari Transaction Abort) ---
		locKey := fmt.Sprintf("%s|%s|%s", countryStr, areaStr, locationStr)
		locationID, exists := locationMap[locKey]
		if !exists {
			err = tx.QueryRow(`
				INSERT INTO locations (country, area, location_name)
				VALUES ($1, $2, $3)
				ON CONFLICT (country, area, location_name) DO UPDATE SET country = EXCLUDED.country
				RETURNING id`, countryStr, areaStr, locationStr).Scan(&locationID)
			if err != nil {
				tx.QueryRow(`SELECT id FROM locations WHERE country=$1 AND area=$2 AND location_name=$3`,
					countryStr, areaStr, locationStr).Scan(&locationID)
			}
			if locationID > 0 {
				locationMap[locKey] = locationID
			}
		}

		// --- 2. Get/Insert Species (Aman dari Transaction Abort) ---
		speciesID, exists := speciesMap[speciesStr]
		if !exists {
			err = tx.QueryRow(`
				INSERT INTO shark_species (species_name, species_group)
				VALUES ($1, 'General')
				ON CONFLICT (species_name) DO UPDATE SET species_name = EXCLUDED.species_name
				RETURNING id`, speciesStr).Scan(&speciesID)
			if err != nil {
				tx.QueryRow(`SELECT id FROM shark_species WHERE species_name=$1`, speciesStr).Scan(&speciesID)
			}
			if speciesID > 0 {
				speciesMap[speciesStr] = speciesID
			}
		}

		// --- 3. Get/Insert Activity (Aman dari Transaction Abort) ---
		activityID, exists := activityMap[activityStr]
		if !exists {
			err = tx.QueryRow(`
				INSERT INTO activities (name, category)
				VALUES ($1, 'Water Activity')
				ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
				RETURNING id`, activityStr).Scan(&activityID)
			if err != nil {
				tx.QueryRow(`SELECT id FROM activities WHERE name=$1`, activityStr).Scan(&activityID)
			}
			if activityID > 0 {
				activityMap[activityStr] = activityID
			}
		}

		// --- 4. Insert Incident ---
		var locPtr, specPtr, actPtr interface{}
		if locationID > 0 {
			locPtr = locationID
		}
		if speciesID > 0 {
			specPtr = speciesID
		}
		if activityID > 0 {
			actPtr = activityID
		}

		_, err = stmtIncident.Exec(dateStr, year, locPtr, specPtr, actPtr, injuryStr, fatalStr)
		if err == nil {
			importedCount++
		}

		// Batasi 30 record agar respon HTTP Vercel instan (< 1.5 detik)
		if importedCount >= 30 {
			break
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return importedCount, nil
}
