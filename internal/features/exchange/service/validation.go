package service

import (
	"fmt"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/features/exchange/gedcom"
)

func validateSnapshot(snapshot gedcom.Snapshot) error {
	for _, person := range snapshot.Persons {
		if person.Input.BirthDate != nil && person.Input.DeathDate != nil && person.Input.DeathDate.Before(*person.Input.BirthDate) {
			return fmt.Errorf("%w: person %s has death date before birth date", ErrInvalid, person.ExternalID)
		}
	}
	children := make(map[string]map[string]struct{})
	for _, relationship := range snapshot.Relationships {
		if relationship.Input.Type != domain.RelationshipParentChild {
			continue
		}
		parent := relationship.Person1ExternalID
		child := relationship.Person2ExternalID
		if relationship.Input.Direction != nil && *relationship.Input.Direction == domain.DirectionChild {
			parent, child = child, parent
		}
		if reaches(child, parent, children, map[string]struct{}{}) {
			return fmt.Errorf("%w: import contains a parent-child cycle", ErrInvalid)
		}
		if children[parent] == nil {
			children[parent] = make(map[string]struct{})
		}
		children[parent][child] = struct{}{}
	}
	return nil
}

func reaches(from, target string, children map[string]map[string]struct{}, visited map[string]struct{}) bool {
	if from == target {
		return true
	}
	if _, ok := visited[from]; ok {
		return false
	}
	visited[from] = struct{}{}
	for child := range children[from] {
		if reaches(child, target, children, visited) {
			return true
		}
	}
	return false
}
