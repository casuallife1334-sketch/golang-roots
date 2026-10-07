package service

import (
	"genealogy-tree/internal/core/domain"
	"strings"
)

func normalizeFamilyInput(input domain.CreateFamilyInput) (domain.CreateFamilyInput, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Metadata == nil {
		input.Metadata = map[string]any{}
	}
	return input, nil
}

func validFamilyMemberRole(role domain.FamilyMemberRole) bool {
	return role == domain.FamilyMemberPartner || role == domain.FamilyMemberParent || role == domain.FamilyMemberChild
}
