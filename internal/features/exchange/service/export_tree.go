package service

import (
	"bytes"
	"context"
	"fmt"
	"genealogy-tree/internal/features/exchange/gedcom"
)

func (s *ExchangeService) Export(ctx context.Context, userID, treeID string) ([]byte, error) {
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return nil, err
	}
	persons, relationships, err := s.exchangeRepository.GetTreeData(ctx, treeID)
	if err != nil {
		return nil, err
	}
	var output bytes.Buffer
	if err := gedcom.Export(persons, relationships, &output); err != nil {
		return nil, fmt.Errorf("export GEDCOM: %w", err)
	}
	return output.Bytes(), nil
}
