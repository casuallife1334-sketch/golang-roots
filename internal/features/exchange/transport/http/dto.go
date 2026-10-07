package http

import "genealogy-tree/internal/features/exchange/gedcom"

type PreviewResponse struct {
	Format        string         `json:"format" example:"gedcom"`
	Persons       int            `json:"persons" example:"12"`
	Relationships int            `json:"relationships" example:"18"`
	Warnings      []gedcom.Issue `json:"warnings,omitempty"`
	Errors        []gedcom.Issue `json:"errors,omitempty"`
}

type ImportResponse struct {
	Format        string `json:"format" example:"gedcom"`
	Persons       int    `json:"persons_imported" example:"12"`
	Relationships int    `json:"relationships_imported" example:"18"`
}
