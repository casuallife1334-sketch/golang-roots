package service

import (
	"context"
	"fmt"
	"io"

	coreerrors "genealogy-tree/internal/core/errors"
)

var ErrPhotoNotFound = fmt.Errorf("%w: photo not found", coreerrors.ErrNotFound)

func (s *PersonsService) GetPersonPhoto(ctx context.Context, userID, treeID, id string) (io.ReadCloser, error) {
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return nil, err
	}
	person, err := s.personsRepository.GetPerson(ctx, treeID, id)
	if err != nil {
		return nil, err
	}
	if person.PhotoURL == nil {
		return nil, ErrPhotoNotFound
	}
	return s.fileStorage.Get(ctx, *person.PhotoURL)
}
