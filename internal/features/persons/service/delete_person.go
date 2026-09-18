package service

import "context"

func (s *PersonsService) DeletePerson(ctx context.Context, userID, treeID, id string) error {
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	person, err := s.personsRepository.GetPerson(ctx, treeID, id)
	if err != nil {
		return err
	}
	if person.PhotoURL != nil {
		if err := s.fileStorage.Delete(ctx, *person.PhotoURL); err != nil {
			return err
		}
	}
	return s.personsRepository.DeletePerson(ctx, treeID, id)
}
