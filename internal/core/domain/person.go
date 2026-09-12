package domain

import "time"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type Person struct {
	ID        string         `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	BirthDate *time.Time     `json:"birth_date,omitempty"`
	DeathDate *time.Time     `json:"death_date,omitempty"`
	Gender    *Gender        `json:"gender,omitempty"`
	PhotoURL  *string        `json:"photo_url,omitempty"`
	Metadata  map[string]any `json:"metadata"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
}

type CreatePersonInput struct {
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	BirthDate *time.Time     `json:"birth_date"`
	DeathDate *time.Time     `json:"death_date"`
	Gender    *Gender        `json:"gender"`
	PhotoURL  *string        `json:"photo_url"`
	Metadata  map[string]any `json:"metadata"`
}

type PatchPersonInput struct {
	FirstName *string         `json:"first_name"`
	LastName  *string         `json:"last_name"`
	BirthDate **time.Time     `json:"birth_date"`
	DeathDate **time.Time     `json:"death_date"`
	Gender    **Gender        `json:"gender"`
	PhotoURL  **string        `json:"photo_url"`
	Metadata  *map[string]any `json:"metadata"`
}
