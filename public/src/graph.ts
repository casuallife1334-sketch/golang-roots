import { Position, type Edge, type Node } from "@xyflow/react";
import type { Person, Relationship } from "./types";
import { formatYears, fullName, initials } from "./utils";

export const CARD_WIDTH = 220;
export const CARD_HEIGHT = 112;
const CARD_GAP = 28;
const GENERATION_GAP = 90;

export type PersonNodeData = {
  person: Person;
  photoUrl?: string;
  selected: boolean;
};
export type JunctionNodeData = { kind: "family-junction" };
export type PersonGraphNode = Node<PersonNodeData, "person">;
export type JunctionGraphNode = Node<JunctionNodeData, "junction">;
export type GraphNode = PersonGraphNode | JunctionGraphNode;

function relationshipEnds(relationship: Relationship) {
  if (relationship.type === "spouse") return null;
  const parentId =
    relationship.direction === "child"
      ? relationship.person2_id
      : relationship.person1_id;
  const childId =
    relationship.direction === "child"
      ? relationship.person1_id
      : relationship.person2_id;
  return { parentId, childId };
}

function getGenerations(people: Person[], relationships: Relationship[]) {
  type Link = { id: string; delta: number };
  const links = new Map<string, Link[]>();
  const parentCount = new Map(people.map((person) => [person.id, 0]));
  const addLink = (from: string, to: string, delta: number) => {
    const current = links.get(from) || [];
    current.push({ id: to, delta });
    links.set(from, current);
  };
  relationships.forEach((relationship) => {
    if (relationship.type === "spouse") {
      addLink(relationship.person1_id, relationship.person2_id, 0);
      addLink(relationship.person2_id, relationship.person1_id, 0);
      return;
    }
    const ends = relationshipEnds(relationship);
    if (!ends) return;
    parentCount.set(ends.childId, (parentCount.get(ends.childId) || 0) + 1);
    addLink(ends.parentId, ends.childId, 1);
    addLink(ends.childId, ends.parentId, -1);
  });
  const generations = new Map<string, number>();
  const visited = new Set<string>();
  people.forEach((person) => {
    if (visited.has(person.id)) return;
    const component: string[] = [];
    const componentQueue = [person.id];
    visited.add(person.id);
    while (componentQueue.length) {
      const current = componentQueue.shift()!;
      component.push(current);
      for (const link of links.get(current) || []) {
        if (!visited.has(link.id)) {
          visited.add(link.id);
          componentQueue.push(link.id);
        }
      }
    }
    const root =
      component.find((id) => parentCount.get(id) === 0) || component[0];
    generations.set(root, 0);
    const queue = [root];
    while (queue.length) {
      const current = queue.shift()!;
      for (const link of links.get(current) || []) {
        if (!generations.has(link.id)) {
          generations.set(link.id, generations.get(current)! + link.delta);
          queue.push(link.id);
        }
      }
    }
  });
  return generations;
}

type Unit = {
  id: string;
  personIds: string[];
  generation: number;
  order: number;
  center?: number;
};

function makeUnits(
  people: Person[],
  relationships: Relationship[],
  generations: Map<string, number>,
) {
  const parent = new Map(people.map((person) => [person.id, person.id]));
  const find = (id: string): string => {
    const current = parent.get(id) || id;
    if (current === id) return id;
    const root = find(current);
    parent.set(id, root);
    return root;
  };
  relationships
    .filter((relationship) => relationship.type === "spouse")
    .forEach((relationship) => {
      const first = find(relationship.person1_id);
      const second = find(relationship.person2_id);
      if (first !== second) parent.set(second, first);
    });
  const units = new Map<string, Unit>();
  people.forEach((person, order) => {
    const id = `family-${find(person.id)}`;
    const unit = units.get(id) || {
      id,
      personIds: [],
      generation: generations.get(person.id) || 0,
      order,
    };
    unit.personIds.push(person.id);
    unit.generation = Math.min(
      unit.generation,
      generations.get(person.id) || 0,
    );
    units.set(id, unit);
  });
  return {
    units: [...units.values()],
    unitByPerson: new Map(
      people.map((person) => [
        person.id,
        units.get(`family-${find(person.id)}`)!,
      ]),
    ),
  };
}

export function buildGraph(
  people: Person[],
  relationships: Relationship[],
  photoUrls: Record<string, string>,
  selectedId?: string,
) {
  const generations = getGenerations(people, relationships);
  const { units, unitByPerson } = makeUnits(people, relationships, generations);
  const personById = new Map(people.map((person) => [person.id, person]));
  const unitParents = new Map<string, Set<string>>();
  const unitChildren = new Map<string, Set<string>>();
  relationships.forEach((relationship) => {
    const ends = relationshipEnds(relationship);
    if (!ends) return;
    const childUnit = unitByPerson.get(ends.childId);
    const parentUnit = unitByPerson.get(ends.parentId);
    if (!childUnit || !parentUnit || childUnit.id === parentUnit.id) return;
    const parents = unitParents.get(childUnit.id) || new Set<string>();
    parents.add(parentUnit.id);
    unitParents.set(childUnit.id, parents);
    const children = unitChildren.get(parentUnit.id) || new Set<string>();
    children.add(childUnit.id);
    unitChildren.set(parentUnit.id, children);
  });
  const byGeneration = new Map<number, Unit[]>();
  units.forEach((unit) =>
    byGeneration.set(unit.generation, [
      ...(byGeneration.get(unit.generation) || []),
      unit,
    ]),
  );
  const unitCenters = new Map<string, number>();
  const unitWidth = (unit: Unit) =>
    unit.personIds.length * CARD_WIDTH + (unit.personIds.length - 1) * CARD_GAP;
  byGeneration.forEach((group) => {
    const total =
      group.reduce((sum, unit) => sum + unitWidth(unit), 0) +
      Math.max(0, group.length - 1) * 90;
    let cursor = -total / 2;
    group
      .sort((a, b) => a.order - b.order)
      .forEach((unit) => {
        unit.center = cursor + unitWidth(unit) / 2;
        unitCenters.set(unit.id, unit.center);
        cursor += unitWidth(unit) + 90;
      });
  });
  const preventOverlaps = (generation: number) => {
    const group = [...(byGeneration.get(generation) || [])].sort(
      (a, b) => (a.center || 0) - (b.center || 0),
    );
    let previousRight = -Infinity;
    group.forEach((unit) => {
      unit.center = Math.max(
        unit.center || 0,
        previousRight + 90 + unitWidth(unit) / 2,
      );
      unitCenters.set(unit.id, unit.center);
      previousRight = unit.center + unitWidth(unit) / 2;
    });
  };
  for (let pass = 0; pass < 8; pass += 1) {
    unitChildren.forEach((children, parentId) => {
      const parentUnit = units.find((unit) => unit.id === parentId);
      const values = [...children]
        .map((id) => unitCenters.get(id))
        .filter((value): value is number => value !== undefined);
      if (!parentUnit || !values.length) return;
      parentUnit.center =
        (parentUnit.center || 0) +
        (values.reduce((sum, value) => sum + value, 0) / values.length -
          (parentUnit.center || 0)) *
          0.45;
      unitCenters.set(parentId, parentUnit.center);
    });
    unitParents.forEach((parents, childId) => {
      const childUnit = units.find((unit) => unit.id === childId);
      const values = [...parents]
        .map((id) => unitCenters.get(id))
        .filter((value): value is number => value !== undefined);
      if (!childUnit || !values.length) return;
      childUnit.center =
        values.reduce((sum, value) => sum + value, 0) / values.length;
      unitCenters.set(childId, childUnit.center);
    });
    byGeneration.forEach((_, generation) => preventOverlaps(generation));
  }
  unitParents.forEach((parents, childId) => {
    const child = units.find((unit) => unit.id === childId);
    const parentUnits = [...parents]
      .map((id) => units.find((unit) => unit.id === id))
      .filter((unit): unit is Unit => Boolean(unit));
    if (!child || parentUnits.length < 2) return;
    const total =
      parentUnits.reduce((sum, unit) => sum + unitWidth(unit), 0) +
      (parentUnits.length - 1) * 90;
    let cursor = (child.center || 0) - total / 2;
    parentUnits.forEach((unit) => {
      unit.center = cursor + unitWidth(unit) / 2;
      unitCenters.set(unit.id, unit.center);
      cursor += unitWidth(unit) + 90;
    });
  });

  const nodes: GraphNode[] = [];
  const personCenters = new Map<string, number>();
  units.forEach((unit) => {
    const width = unitWidth(unit);
    const start = (unit.center || 0) - width / 2;
    unit.personIds.forEach((personId, index) => {
      const x = start + index * (CARD_WIDTH + CARD_GAP);
      const person = personById.get(personId)!;
      personCenters.set(personId, x + CARD_WIDTH / 2);
      nodes.push({
        id: personId,
        type: "person",
        className: "person-flow-node",
        position: { x, y: unit.generation * (CARD_HEIGHT + GENERATION_GAP) },
        data: {
          person,
          photoUrl: photoUrls[personId],
          selected: personId === selectedId,
        },
        sourcePosition: Position.Bottom,
        targetPosition: Position.Top,
      });
    });
  });
  const edges: Edge[] = relationships
    .filter((relationship) => relationship.type === "spouse")
    .map((relationship) => ({
      id: `spouse-${relationship.id}`,
      source: relationship.person1_id,
      target: relationship.person2_id,
      sourceHandle: "spouse-out",
      targetHandle: "spouse-in",
      type: "straight",
      style: { stroke: "#7c8791", strokeWidth: 2.5 },
    }));
  const parentSets = new Map<string, Set<string>>();
  relationships.forEach((relationship) => {
    const ends = relationshipEnds(relationship);
    if (
      !ends ||
      !personById.has(ends.parentId) ||
      !personById.has(ends.childId)
    )
      return;
    const parents = parentSets.get(ends.childId) || new Set<string>();
    parents.add(ends.parentId);
    parentSets.set(ends.childId, parents);
  });
  const junctions = new Map<string, string>();
  const parentLinks = new Map<
    string,
    { parentId: string; childId: string }[]
  >();
  parentSets.forEach((parents, childId) => {
    const parentIds = [...parents].sort();
    const key = `${parentIds.join("|")}::${generations.get(childId)}`;
    const junctionId = junctions.get(key) || `junction-${junctions.size}`;
    junctions.set(key, junctionId);
    const links = parentIds.map((parentId) => ({ parentId, childId }));
    parentLinks.set(junctionId, [
      ...(parentLinks.get(junctionId) || []),
      ...links,
    ]);
  });
  junctions.forEach((junctionId, key) => {
    const [parentKey, generationValue] = key.split("::");
    const parentIds = parentKey.split("|");
    const x =
      parentIds.reduce((sum, id) => sum + (personCenters.get(id) || 0), 0) /
      Math.max(1, parentIds.length);
    const y =
      Number(generationValue) * (CARD_HEIGHT + GENERATION_GAP) -
      GENERATION_GAP / 2 -
      1;
    nodes.push({
      id: junctionId,
      type: "junction",
      className: "junction-flow-node",
      position: { x: x - 0.5, y },
      data: { kind: "family-junction" },
      sourcePosition: Position.Bottom,
      targetPosition: Position.Top,
      selectable: false,
      draggable: false,
    });
  });
  parentLinks.forEach((links, junctionId) => {
    const unique = [
      ...new Map(
        links.map((link) => [`${link.parentId}-${link.childId}`, link]),
      ).values(),
    ];
    [...new Set(unique.map((link) => link.parentId))].forEach((parentId) =>
      edges.push({
        id: `parent-${junctionId}-${parentId}`,
        source: parentId,
        target: junctionId,
        type: "smoothstep",
        style: { stroke: "#a0a5aa", strokeWidth: 1.8 },
      }),
    );
    [...new Set(unique.map((link) => link.childId))].forEach((childId) =>
      edges.push({
        id: `child-${junctionId}-${childId}`,
        source: junctionId,
        target: childId,
        type: "smoothstep",
        style: { stroke: "#a0a5aa", strokeWidth: 1.8 },
      }),
    );
  });
  return { nodes, edges };
}

export function personNodeSubtitle(person: Person) {
  return [formatYears(person), person.metadata?.city]
    .filter(Boolean)
    .join(" · ");
}
export { fullName, initials };
