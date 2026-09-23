import dagre from "@dagrejs/dagre";
import type { FamilyModel } from "./model";
export const CARD_WIDTH = 240,
  CARD_HEIGHT = 132;
export type Point = { x: number; y: number };
// A person has one card. Only non-overlapping pairs are packed together;
// remarriages remain distinct relations rather than a transitive spouse union.
export function layoutFamilies(model: FamilyModel, compact = false) {
  const gap = compact ? 32 : 48;
  const groups = new Map<string, string[]>(),
    groupOf = new Map<string, string>();
  const canPair = (a: string, b: string) => {
    const visited = new Set<string>(),
      queue = [a];
    while (queue.length) {
      const id = queue.pop()!;
      if (id === b) return false;
      if (visited.has(id)) continue;
      visited.add(id);
      queue.push(...(model.children.get(id) ?? []));
    }
    return true;
  };
  const pairs = [
    ...model.spouses,
    ...model.families
      .filter((family) => family.parents.length === 2)
      .map((family) => family.parents as [string, string]),
  ];
  for (const [a, b] of pairs) {
    if (groupOf.has(a) || groupOf.has(b) || !canPair(a, b) || !canPair(b, a))
      continue;
    const id = `group:${a}`;
    groups.set(id, [a, b]);
    groupOf.set(a, id);
    groupOf.set(b, id);
  }
  for (const person of model.persons)
    if (!groupOf.has(person.id)) {
      const id = `group:${person.id}`;
      groups.set(id, [person.id]);
      groupOf.set(person.id, id);
    }
  const graph = new dagre.graphlib.Graph()
    .setGraph({
      rankdir: "TB",
      nodesep: compact ? 48 : 80,
      ranksep: compact ? 64 : 88,
      marginx: 24,
      marginy: 24,
    })
    .setDefaultEdgeLabel(() => ({}));
  for (const [id, members] of groups)
    graph.setNode(id, {
      width: members.length * CARD_WIDTH + (members.length - 1) * gap,
      height: CARD_HEIGHT,
    });
  for (const family of model.families)
    for (const parent of family.parents)
      for (const child of family.children) {
        const a = groupOf.get(parent)!,
          b = groupOf.get(child)!;
        if (a !== b) graph.setEdge(a, b);
      }
  dagre.layout(graph);
  const positions = new Map<string, Point>();
  for (const [id, members] of groups) {
    const node = graph.node(id);
    const parentCenter = (person: string) => {
      const parentGroups = [
        ...new Set(
          [...(model.parents.get(person) ?? [])].map(
            (parent) => groupOf.get(parent)!,
          ),
        ),
      ];
      return parentGroups.length
        ? parentGroups.reduce((sum, parent) => sum + graph.node(parent).x, 0) /
            parentGroups.length
        : node.x;
    };
    [...members]
      .sort((a, b) => parentCenter(a) - parentCenter(b) || a.localeCompare(b))
      .forEach((person, index) =>
        positions.set(person, {
          x: node.x - node.width / 2 + index * (CARD_WIDTH + gap),
          y: node.y - CARD_HEIGHT / 2,
        }),
      );
  }
  return positions;
}
