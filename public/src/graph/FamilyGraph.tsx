import { Background, Controls, MiniMap, ReactFlow } from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { useMemo, useState } from "react";
import { Viewport } from "./Viewport";
import { buildGraph, type PersonGraphNode } from "../graph";
import type { Person, Relationship } from "../types";
import { PersonNode } from "./PersonNode";
import { FamilyEdge } from "./FamilyEdge";
import { usePreferences } from "../data/preferences";
import { Notice } from "../shared/ui";
import "./graph.css";
const nodeTypes = { person: PersonNode },
  edgeTypes = { family: FamilyEdge };
export function FamilyGraph({
  treeId,
  people,
  relationships,
  selectedId,
  onSelect,
}: {
  treeId: string;
  people: Person[];
  relationships: Relationship[];
  selectedId?: string;
  onSelect: (person: Person) => void;
}) {
  const prefs = usePreferences();
  const [container, setContainer] = useState<HTMLDivElement | null>(null);
  const graph = useMemo(
    () => buildGraph(people, relationships),
    [people, relationships],
  );
  const nodes: PersonGraphNode[] = useMemo(
    () =>
      graph.nodes.map((node) => ({
        ...node,
        data: {
          ...node.data,
          treeId,
          selected: node.id === selectedId,
          portraits: prefs.portraits,
        },
      })),
    [graph.nodes, treeId, selectedId, prefs.portraits],
  );
  return (
    <div ref={setContainer} style={{ height: "100%", width: "100%" }}>
      {graph.warnings.length > 0 && (
        <div className="graph-warning">
          <Notice>{[...new Set(graph.warnings)].join(". ")}</Notice>
        </div>
      )}
      <ReactFlow
        nodes={nodes}
        edges={graph.edges}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        fitView
        fitViewOptions={{ padding: 0.2, maxZoom: 1 }}
        minZoom={0.1}
        maxZoom={2}
        proOptions={{ hideAttribution: true }}
        nodesConnectable={false}
        nodesDraggable={false}
        onlyRenderVisibleElements
        onNodeClick={(_, node) => onSelect(node.data.person)}
      >
        <Viewport
          container={container}
          topology={JSON.stringify([
            graph.nodes.map((node) => [node.id, node.position]),
            graph.edges.map((edge) => edge.id),
          ])}
        />
        <Background color="#d5d2cf" gap={24} size={1} />
        <Controls showInteractive={false} />
        {prefs.minimap && <MiniMap nodeColor="#a1adb5" />}
      </ReactFlow>
    </div>
  );
}
