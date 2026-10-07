package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/features/exchange/gedcom"
	"io"
)

func (s *ExchangeService) Import(ctx context.Context, userID, treeID string, input io.Reader) (ImportResult, gedcom.Snapshot, error) {
	snapshot, err := gedcom.Parse(input)
	if err != nil {
		return ImportResult{}, snapshot, err
	}
	if len(snapshot.Errors) > 0 {
		return ImportResult{}, snapshot, fmt.Errorf("%w: GEDCOM preview contains validation errors", ErrInvalid)
	}
	if err := validateSnapshot(snapshot); err != nil {
		return ImportResult{}, snapshot, err
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return ImportResult{}, snapshot, err
	}
	persons, relationships, err := s.exchangeRepository.ImportSnapshot(ctx, treeID, snapshot)
	if err != nil {
		return ImportResult{}, snapshot, err
	}
	return ImportResult{Persons: persons, Relationships: relationships}, snapshot, nil
}
