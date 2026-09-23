import { api } from "../../api";
import type { Person, PersonInput } from "../../types";

export type PersonForm = {
  first_name: string;
  last_name: string;
  patronymic: string;
  birth_date: string;
  death_date: string;
  gender: string;
  city: string;
  occupation: string;
};
export function personInput(
  form: PersonForm,
  person: Person | null,
): PersonInput {
  return {
    first_name: form.first_name.trim(),
    last_name: form.last_name.trim(),
    birth_date: form.birth_date ? `${form.birth_date}T00:00:00Z` : null,
    death_date: form.death_date ? `${form.death_date}T00:00:00Z` : null,
    gender: (form.gender || null) as PersonInput["gender"],
    metadata: {
      ...person?.metadata,
      patronymic: form.patronymic.trim(),
      city: form.city.trim(),
      occupation: form.occupation.trim(),
    },
  };
}
// Retain the successful create result even when the following photo request fails.
export class PersonSave {
  saved: Person | null;
  constructor(person: Person | null) {
    this.saved = person;
  }
  async run(
    tree: string,
    body: PersonInput,
    photo: File | undefined,
    removePhoto: boolean,
    onSaved: (person: Person) => void,
  ) {
    this.saved = this.saved
      ? await api.patchPerson(tree, this.saved.id, body)
      : await api.createPerson(tree, body);
    onSaved(this.saved);
    try {
      if (photo) this.saved = await api.uploadPhoto(tree, this.saved.id, photo);
      else if (removePhoto) {
        await api.deletePhoto(tree, this.saved.id);
        this.saved = { ...this.saved, photo_url: null };
      }
      onSaved(this.saved);
      return this.saved;
    } catch {
      throw new Error(
        "Человек сохранён, но фотографию изменить не удалось. Повторите сохранение — дубликат создан не будет.",
      );
    }
  }
}
