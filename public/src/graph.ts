import type { Node } from "@xyflow/react";
import type { Person, Relationship } from "./types";
import { familyModel } from "./graph/model";
import { layoutFamilies, CARD_WIDTH, CARD_HEIGHT } from "./graph/layout";
import { buildEdges } from "./graph/edges";
export { CARD_WIDTH, CARD_HEIGHT } from "./graph/layout";
export type PersonGraphNode = Node<
  { person: Person; treeId: string; selected: boolean; portraits: boolean },
  "person"
>;
export function buildGraph(
  people: Person[],
  relationships: Relationship[],
  compact = false,
) {
  const model = familyModel(people, relationships),
    positions = layoutFamilies(model, compact);
  return {
    nodes: model.persons.map((person) => ({
      id: person.id,
      type: "person" as const,
      className: "person-flow-node",
      position: positions.get(person.id)!,
      width: CARD_WIDTH,
      height: CARD_HEIGHT,
      data: { person, treeId: "", selected: false, portraits: true },
    })),
    edges: buildEdges(model, positions),
    warnings: model.warnings,
  };
}
