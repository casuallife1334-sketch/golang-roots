package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/features/exchange/gedcom"
)

var ErrInvalid = fmt.Errorf("%w: invalid GEDCOM import", coreerrors.ErrInvalidArgument)
var ErrTooLarge = fmt.Errorf("%w: GEDCOM file is too large", coreerrors.ErrPayloadTooLarge)

type ExchangeRepository interface {
	GetTreeData(context.Context, string) ([]domain.Person, []domain.Relationship, error)
	ImportSnapshot(context.Context, string, gedcom.Snapshot) (int, int, error)
}

type TreeAccess interface {
	CanReadTree(context.Context, string, string) error
	CanWriteTree(context.Context, string, string) error
}

type ImportResult struct {
	Persons       int
	Relationships int
}

type ExchangeService struct {
	exchangeRepository ExchangeRepository
	treeAccess         TreeAccess
}

func NewExchangeService(exchangeRepository ExchangeRepository, treeAccess TreeAccess) *ExchangeService {
	return &ExchangeService{
		exchangeRepository: exchangeRepository,
		treeAccess:         treeAccess,
	}
}
