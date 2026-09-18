package service

import (
	"context"
	"fmt"
	"io"

	"genealogy-tree/internal/core/domain"
)

func (s *PersonsService) UploadPersonPhoto(ctx context.Context, userID, treeID, id, contentType string, file io.Reader) (domain.Person, error) {
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Person{}, err
	}
	person, err := s.personsRepository.GetPerson(ctx, treeID, id)
	if err != nil {
		return domain.Person{}, err
	}

	key := "persons/" + id + "/photo"
	photoURL, err := s.fileStorage.Put(ctx, key, file, contentType)
	if err != nil {
		return domain.Person{}, fmt.Errorf("put person photo: %w", err)
	}

	updatedPerson, err := s.personsRepository.UpdatePersonPhoto(ctx, treeID, id, &photoURL)
	if err != nil {
		_ = s.fileStorage.Delete(ctx, photoURL)
		return domain.Person{}, err
	}
	if person.PhotoURL != nil && *person.PhotoURL != photoURL {
		_ = s.fileStorage.Delete(ctx, *person.PhotoURL)
	}
	return updatedPerson, nil
}
