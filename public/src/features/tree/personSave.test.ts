import { afterEach, expect, it, vi } from "vitest";
import { api } from "../../api";
import { personInput, PersonSave, type PersonForm } from "./personSave";
const person = {
  id: "p",
  first_name: "A",
  last_name: "B",
  metadata: { city: "Old", patronymic: "Old", occupation: "Old" },
  created_at: "",
};
const form: PersonForm = {
  first_name: "A",
  last_name: "B",
  city: "",
  patronymic: "",
  occupation: "",
  birth_date: "",
  death_date: "",
  gender: "",
};
afterEach(() => vi.restoreAllMocks());
it("clears metadata explicitly", () => {
  expect(personInput(form, person).metadata).toEqual({
    city: "",
    patronymic: "",
    occupation: "",
  });
});
it("does not create another person when retrying a failed photo upload", async () => {
  const create = vi.spyOn(api, "createPerson").mockResolvedValue(person);
  vi.spyOn(api, "patchPerson").mockResolvedValue(person);
  vi.spyOn(api, "uploadPhoto")
    .mockRejectedValueOnce(new Error("offline"))
    .mockResolvedValueOnce({ ...person, photo_url: "/photo" });
  const save = new PersonSave(null),
    onSaved = vi.fn(),
    file = {} as File;
  await expect(
    save.run("tree", personInput(form, null), file, false, onSaved),
  ).rejects.toThrow("Человек сохранён");
  expect(onSaved).toHaveBeenCalledWith(person);
  await save.run("tree", personInput(form, null), file, false, onSaved);
  expect(create).toHaveBeenCalledTimes(1);
});
