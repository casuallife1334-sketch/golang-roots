package service

import (
	"genealogy-tree/internal/features/exchange/gedcom"
	"io"
)

func (s *ExchangeService) Preview(input io.Reader) (gedcom.Snapshot, error) {
	snapshot, err := gedcom.Parse(input)
	if err != nil {
		return snapshot, err
	}
	if err := validateSnapshot(snapshot); err != nil {
		snapshot.Errors = append(snapshot.Errors, gedcom.Issue{Message: err.Error()})
	}
	return snapshot, nil
}
