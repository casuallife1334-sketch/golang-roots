package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

type FamilyRepository interface {
	CreateFamily(context.Context, string, domain.CreateFamilyInput) (domain.Family, error)
	GetFamilies(context.Context, string) ([]domain.Family, error)
	GetFamily(context.Context, string, string) (domain.Family, error)
	PatchFamily(context.Context, string, string, domain.PatchFamilyInput) (domain.Family, error)
	DeleteFamily(context.Context, string, string) error
}
