import type { Family, Person } from "../../types";
import { fullName } from "../../utils";

export function familyOptionLabel(family: Family, people: Person[]): string {
  const peopleById = new Map(people.map((person) => [person.id, person]));
  const adultIds = new Set(
    family.members
      .filter((member) => member.role !== "child")
      .map((member) => member.person_id),
  );
  const names = [...adultIds]
    .map((id) => peopleById.get(id))
    .filter((person): person is Person => Boolean(person))
    .map(fullName);

  return family.name || names.join(" и ") || family.id;
}
