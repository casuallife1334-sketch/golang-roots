package service

import (
	"context"
	"fmt"
	"io"

	coreerrors "genealogy-tree/internal/core/errors"
)

var ErrPhotoNotFound = fmt.Errorf("%w: photo not found", coreerrors.ErrNotFound)

func (s *PersonsService) GetPersonPhoto(ctx context.Context, id string) (io.ReadCloser, error) {
	person, err := s.personsRepository.GetPerson(ctx, id)
	if err != nil {
		return nil, err
	}
	if person.PhotoURL == nil {
		return nil, ErrPhotoNotFound
	}
	return s.fileStorage.Get(ctx, *person.PhotoURL)
}
