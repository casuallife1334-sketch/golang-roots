package service

import (
	"context"
	"fmt"
)

func (s *PersonsService) DeletePersonPhoto(ctx context.Context, userID, treeID, id string) error {
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	person, err := s.personsRepository.GetPerson(ctx, treeID, id)
	if err != nil {
		return err
	}
	if person.PhotoURL == nil {
		return ErrPhotoNotFound
	}
	if err := s.fileStorage.Delete(ctx, *person.PhotoURL); err != nil {
		return fmt.Errorf("delete person photo: %w", err)
	}
	_, err = s.personsRepository.UpdatePersonPhoto(ctx, treeID, id, nil)
	return err
}
