package service

import (
	"context"
	"fmt"
	"io"

	"genealogy-tree/internal/core/domain"
)

func (s *PersonsService) UploadPersonPhoto(ctx context.Context, id, contentType string, file io.Reader) (domain.Person, error) {
	person, err := s.personsRepository.GetPerson(ctx, id)
	if err != nil {
		return domain.Person{}, err
	}

	key := "persons/" + id + "/photo"
	photoURL, err := s.fileStorage.Put(ctx, key, file, contentType)
	if err != nil {
		return domain.Person{}, fmt.Errorf("put person photo: %w", err)
	}

	updatedPerson, err := s.personsRepository.UpdatePersonPhoto(ctx, id, &photoURL)
	if err != nil {
		_ = s.fileStorage.Delete(ctx, photoURL)
		return domain.Person{}, err
	}
	if person.PhotoURL != nil && *person.PhotoURL != photoURL {
		_ = s.fileStorage.Delete(ctx, *person.PhotoURL)
	}
	return updatedPerson, nil
}
