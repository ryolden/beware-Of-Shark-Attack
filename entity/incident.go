package entity

import "time"

type Location struct {
	ID           int       `json:"id"`
	Country      string    `json:"country"`
	Area         string    `json:"area"`
	LocationName string    `json:"location_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type SharkSpecies struct {
	ID           int       `json:"id"`
	SpeciesName  string    `json:"species_name"`
	SpeciesGroup string    `json:"species_group"`
	CreatedAt    time.Time `json:"created_at"`
}

type Incident struct {
	ID           int           `json:"id"`
	IncidentDate string        `json:"incident_date"`
	Year         int           `json:"year"`
	UserID       *int          `json:"user_id,omitempty"`
	LocationID   int           `json:"location_id"`
	SpeciesID    int           `json:"species_id"`
	ActivityID   int           `json:"activity_id"`
	Injury       string        `json:"injury"`
	Fatal        string        `json:"fatal"`
	CreatedAt    time.Time     `json:"created_at"`
	Location     *Location     `json:"location,omitempty"`
	Species      *SharkSpecies `json:"species,omitempty"`
	Activity     *Activity     `json:"activity,omitempty"`
}
