package service

import (
	"context"
	"fmt"
)

func (s *PersonsService) DeletePersonPhoto(ctx context.Context, id string) error {
	person, err := s.personsRepository.GetPerson(ctx, id)
	if err != nil {
		return err
	}
	if person.PhotoURL == nil {
		return ErrPhotoNotFound
	}
	if err := s.fileStorage.Delete(ctx, *person.PhotoURL); err != nil {
		return fmt.Errorf("delete person photo: %w", err)
	}
	_, err = s.personsRepository.UpdatePersonPhoto(ctx, id, nil)
	return err
}
