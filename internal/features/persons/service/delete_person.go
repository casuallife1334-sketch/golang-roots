package service

import "context"

func (s *PersonsService) DeletePerson(ctx context.Context, id string) error {
	person, err := s.personsRepository.GetPerson(ctx, id)
	if err != nil {
		return err
	}
	if person.PhotoURL != nil {
		if err := s.fileStorage.Delete(ctx, *person.PhotoURL); err != nil {
			return err
		}
	}
	return s.personsRepository.DeletePerson(ctx, id)
}
