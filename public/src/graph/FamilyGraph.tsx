import { Background, Controls, MiniMap, ReactFlow } from "@xyflow/react";
import { useMemo } from "react";
import { buildGraph } from "../graph";
import type { Person, Relationship } from "../types";
import { JunctionNode } from "./JunctionNode";
import { PersonNode } from "./PersonNode";

const nodeTypes = { person: PersonNode, junction: JunctionNode };
type Props = {
  people: Person[];
  relationships: Relationship[];
  photoUrls: Record<string, string>;
  selectedId?: string;
  onSelect: (person: Person) => void;
};

export function FamilyGraph({
  people,
  relationships,
  photoUrls,
  selectedId,
  onSelect,
}: Props) {
  const graph = useMemo(
    () => buildGraph(people, relationships, photoUrls, selectedId),
    [people, relationships, photoUrls, selectedId],
  );
  return (
    <ReactFlow
      nodes={graph.nodes}
      edges={graph.edges}
      nodeTypes={nodeTypes}
      fitView
      fitViewOptions={{ padding: 0.25 }}
      nodesConnectable={false}
      nodesDraggable={false}
      onNodeClick={(_, node) => {
        const person = people.find((item) => item.id === node.id);
        if (person) onSelect(person);
      }}
      proOptions={{ hideAttribution: true }}
    >
      <Background color="#d5d2cf" gap={24} size={1} />
      <Controls showInteractive={false} />
      <MiniMap
        nodeColor={(node) => (node.data.selected ? "#3478ff" : "#9aa0a6")}
      />
    </ReactFlow>
  );
}
