import type { Person, Relationship, RelationshipInput } from "../types";
export function parentChild(relation: RelationshipInput | Relationship) {
  if (relation.type === "spouse") return null;
  return relation.direction === "child"
    ? { parentId: relation.person2_id, childId: relation.person1_id }
    : { parentId: relation.person1_id, childId: relation.person2_id };
}
export function reaches(
  from: string,
  target: string,
  children: Map<string, Set<string>>,
) {
  const stack = [from],
    visited = new Set<string>();
  while (stack.length) {
    const id = stack.pop()!;
    if (id === target) return true;
    if (visited.has(id)) continue;
    visited.add(id);
    stack.push(...(children.get(id) ?? []));
  }
  return false;
}
export function familyModel(people: Person[], relationships: Relationship[]) {
  const persons = [...people].sort((a, b) => a.id.localeCompare(b.id));
  const ids = new Set(persons.map((person) => person.id));
  const children = new Map<string, Set<string>>(),
    parents = new Map<string, Set<string>>();
  const spouses: [string, string][] = [];
  const warnings: string[] = [];
  const seen = new Set<string>();
  for (const relation of [...relationships].sort((a, b) =>
    a.id.localeCompare(b.id),
  )) {
    if (
      !ids.has(relation.person1_id) ||
      !ids.has(relation.person2_id) ||
      relation.person1_id === relation.person2_id
    ) {
      warnings.push(
        "Пропущена связь с отсутствующим участником или с самим собой",
      );
      continue;
    }
    const link = parentChild(relation);
    if (!link) {
      const pair = [relation.person1_id, relation.person2_id].sort() as [
        string,
        string,
      ];
      const key = JSON.stringify(pair);
      if (!seen.has(key)) spouses.push(pair);
      seen.add(key);
      continue;
    }
    if (reaches(link.childId, link.parentId, children)) {
      warnings.push(
        "Обнаружен цикл родительских связей. Замыкающая связь не отображена",
      );
      continue;
    }
    if (!children.has(link.parentId)) children.set(link.parentId, new Set());
    if (!parents.has(link.childId)) parents.set(link.childId, new Set());
    children.get(link.parentId)!.add(link.childId);
    parents.get(link.childId)!.add(link.parentId);
  }
  const families = new Map<
    string,
    { id: string; parents: string[]; children: string[] }
  >();
  for (const [child, parentSet] of [...parents].sort(([a], [b]) =>
    a.localeCompare(b),
  )) {
    const parentIds = [...parentSet].sort(),
      key = JSON.stringify(parentIds);
    const family = families.get(key) ?? {
      id: key,
      parents: parentIds,
      children: [],
    };
    family.children.push(child);
    families.set(key, family);
  }
  return {
    persons,
    spouses,
    families: [...families.values()],
    parents,
    children,
    warnings,
  };
}
export type FamilyModel = ReturnType<typeof familyModel>;
export function relationshipProblem(
  people: Person[],
  relationships: Relationship[],
  next: RelationshipInput,
): string | null {
  const ids = new Set(people.map((person) => person.id));
  if (
    !ids.has(next.person1_id) ||
    !ids.has(next.person2_id) ||
    next.person1_id === next.person2_id
  )
    return "Выберите двух разных людей";
  const model = familyModel(people, relationships),
    link = parentChild(next);
  if (!link)
    return model.spouses.some(
      (pair) =>
        pair.includes(next.person1_id) && pair.includes(next.person2_id),
    )
      ? "Такая связь уже существует"
      : null;
  if (model.children.get(link.parentId)?.has(link.childId))
    return "Такая связь уже существует";
  if (reaches(link.childId, link.parentId, model.children))
    return "Связь создаёт цикл: человек не может быть собственным предком";
  return null;
}
