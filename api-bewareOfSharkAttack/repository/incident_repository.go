package repository

import (
	"api-bewareOfSharkAttack/entity"
	"database/sql"
	"errors"
)

type IncidentRepository struct {
	DB *sql.DB
}

func NewIncidentRepository(db *sql.DB) *IncidentRepository {
	return &IncidentRepository{DB: db}
}

func (r *IncidentRepository) GetAll(limit, offset int) ([]entity.Incident, error) {
	query := `
		SELECT 
			i.id, i.incident_date, i.year, i.injury, i.fatal, i.created_at,
			l.id, l.country, l.area, l.location_name,
			s.id, s.species_name, s.species_group,
			a.id, a.name, a.category
		FROM incidents i
		LEFT JOIN locations l ON i.location_id = l.id
		LEFT JOIN shark_species s ON i.species_id = s.id
		LEFT JOIN activities a ON i.activity_id = a.id
		ORDER BY i.id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []entity.Incident
	for rows.Next() {
		var inc entity.Incident
		var loc entity.Location
		var spc entity.SharkSpecies
		var act entity.Activity

		err := rows.Scan(
			&inc.ID, &inc.IncidentDate, &inc.Year, &inc.Injury, &inc.Fatal, &inc.CreatedAt,
			&loc.ID, &loc.Country, &loc.Area, &loc.LocationName,
			&spc.ID, &spc.SpeciesName, &spc.SpeciesGroup,
			&act.ID, &act.Name, &act.Category,
		)
		if err != nil {
			return nil, err
		}

		inc.Location = &loc
		inc.Species = &spc
		inc.Activity = &act
		incidents = append(incidents, inc)
	}

	return incidents, nil
}

func (r *IncidentRepository) GetByID(id int) (*entity.Incident, error) {
	query := `
		SELECT 
			i.id, i.incident_date, i.year, i.injury, i.fatal, i.created_at,
			l.id, l.country, l.area, l.location_name,
			s.id, s.species_name, s.species_group,
			a.id, a.name, a.category
		FROM incidents i
		LEFT JOIN locations l ON i.location_id = l.id
		LEFT JOIN shark_species s ON i.species_id = s.id
		LEFT JOIN activities a ON i.activity_id = a.id
		WHERE i.id = $1
	`
	var inc entity.Incident
	var loc entity.Location
	var spc entity.SharkSpecies
	var act entity.Activity

	err := r.DB.QueryRow(query, id).Scan(
		&inc.ID, &inc.IncidentDate, &inc.Year, &inc.Injury, &inc.Fatal, &inc.CreatedAt,
		&loc.ID, &loc.Country, &loc.Area, &loc.LocationName,
		&spc.ID, &spc.SpeciesName, &spc.SpeciesGroup,
		&act.ID, &act.Name, &act.Category,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("data insiden tidak ditemukan")
	} else if err != nil {
		return nil, err
	}

	inc.Location = &loc
	inc.Species = &spc
	inc.Activity = &act
	return &inc, nil
}

func (r *IncidentRepository) Create(inc *entity.Incident) error {
	query := `
		INSERT INTO incidents (incident_date, year, user_id, location_id, species_id, activity_id, injury, fatal)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	return r.DB.QueryRow(query,
		inc.IncidentDate, inc.Year, inc.UserID, inc.LocationID,
		inc.SpeciesID, inc.ActivityID, inc.Injury, inc.Fatal,
	).Scan(&inc.ID, &inc.CreatedAt)
}

func (r *IncidentRepository) Delete(id int) error {
	query := `DELETE FROM incidents WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data insiden tidak ditemukan")
	}
	return nil
}
