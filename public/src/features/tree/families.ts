import { familyModel } from "../../graph/model";
import type { Person, Relationship } from "../../types";
export interface Family {
  id: string;
  members: Person[];
  children: Person[];
  name: string;
}
export function deriveFamilies(
  people: Person[],
  relationships: Relationship[],
): Family[] {
  const model = familyModel(people, relationships),
    map = new Map(people.map((person) => [person.id, person]));
  const groups = new Map<string, Family>();
  const name = (members: Person[]) => {
    const surnames = [...new Set(members.map((person) => person.last_name))];
    return surnames.length === 1
      ? `Семья ${surnames[0]}`
      : members.map((person) => person.first_name).join(" и ");
  };
  for (const pair of model.spouses) {
    const members = pair.map((id) => map.get(id)!);
    const id = JSON.stringify([...pair].sort());
    groups.set(id, { id, members, children: [], name: name(members) });
  }
  for (const family of model.families) {
    const members = family.parents.map((id) => map.get(id)!);
    groups.set(family.id, {
      id: family.id,
      members,
      children: family.children.map((id) => map.get(id)!),
      name: name(members),
    });
  }
  return [...groups.values()].sort((a, b) => a.name.localeCompare(b.name));
}
