package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

type FamilyMembershipRepository interface {
	AddMember(context.Context, string, string, domain.AddFamilyMemberInput) error
	RemoveMember(context.Context, string, string, string) error
}
