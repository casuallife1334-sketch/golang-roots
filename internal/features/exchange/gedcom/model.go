package gedcom

import (
	"genealogy-tree/internal/core/domain"
	"time"
)

type Issue struct {
	Line    int    `json:"line,omitempty"`
	Message string `json:"message"`
}

type PersonRecord struct {
	ExternalID string
	Input      domain.CreatePersonInput
}

type RelationshipRecord struct {
	Person1ExternalID string
	Person2ExternalID string
	Input             domain.CreateRelationshipInput
}

type Snapshot struct {
	Persons       []PersonRecord
	Relationships []RelationshipRecord
	Warnings      []Issue
	Errors        []Issue
}

func dateValue(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02")
}
