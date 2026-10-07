package service

import "context"

type FamilyRelationshipRepository interface {
	AttachRelationship(context.Context, string, string, string) error
	DetachRelationship(context.Context, string, string, string) error
}
