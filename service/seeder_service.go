package service

import (
	"database/sql"
	"encoding/csv"
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

	importedCount := 0
	maxRows := 50 // Dibatasi agar tidak hit timeout Vercel Serverless

	for {
		if importedCount >= maxRows {
			break
		}

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

		// 1. Get or Insert Location
		var locationID int
		err = tx.QueryRow(`SELECT id FROM locations WHERE country=$1 AND area=$2 AND location_name=$3`,
			countryStr, areaStr, locationStr).Scan(&locationID)
		if err == sql.ErrNoRows {
			_ = tx.QueryRow(`
				INSERT INTO locations (country, area, location_name)
				VALUES ($1, $2, $3)
				RETURNING id`, countryStr, areaStr, locationStr).Scan(&locationID)
		}

		// 2. Get or Insert Species
		var speciesID int
		err = tx.QueryRow(`SELECT id FROM shark_species WHERE species_name=$1`, speciesStr).Scan(&speciesID)
		if err == sql.ErrNoRows {
			_ = tx.QueryRow(`
				INSERT INTO shark_species (species_name, species_group)
				VALUES ($1, 'General')
				RETURNING id`, speciesStr).Scan(&speciesID)
		}

		// 3. Get or Insert Activity
		var activityID int
		err = tx.QueryRow(`SELECT id FROM activities WHERE name=$1`, activityStr).Scan(&activityID)
		if err == sql.ErrNoRows {
			_ = tx.QueryRow(`
				INSERT INTO activities (name, category)
				VALUES ($1, 'Water Activity')
				RETURNING id`, activityStr).Scan(&activityID)
		}

		// 4. Insert Incident
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

		_, err = tx.Exec(`
			INSERT INTO incidents (incident_date, year, location_id, species_id, activity_id, injury, fatal)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			dateStr, year, locPtr, specPtr, actPtr, injuryStr, fatalStr)

		if err == nil {
			importedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return importedCount, nil
}
