import type { Edge } from "@xyflow/react";
import type { FamilyModel } from "./model";
import { CARD_HEIGHT, CARD_WIDTH, type Point } from "./layout";
export function buildEdges(model: FamilyModel, positions: Map<string, Point>) {
  const edges: Edge[] = [];
  model.spouses.forEach(([a, b]) => {
    const first = positions.get(a)!,
      second = positions.get(b)!;
    const left = first.x <= second.x ? a : b,
      right = left === a ? b : a;
    edges.push({
      id: `spouse:${JSON.stringify([a, b])}`,
      source: left,
      target: right,
      sourceHandle: "spouse-out",
      targetHandle: "spouse-in",
      type: first.y === second.y ? "straight" : "smoothstep",
      style: { stroke: "#7c8791", strokeWidth: 2 },
    });
  });
  model.families.forEach((family, index) => {
    const parents = family.parents.map((id) => positions.get(id)!);
    const children = family.children.map((id) => positions.get(id)!);
    const top = Math.max(...parents.map((point) => point.y + CARD_HEIGHT));
    const bottom = Math.min(...children.map((point) => point.y));
    const center =
      (Math.min(...parents.map((point) => point.x)) +
        Math.max(...parents.map((point) => point.x + CARD_WIDTH))) /
      2;
    const lane =
      top + Math.max(16, (bottom - top) * (0.35 + (index % 3) * 0.12));
    // One SVG path per real parent set: a shared bus cannot connect in-laws.
    const xs = [...parents, ...children].map(
      (point) => point.x + CARD_WIDTH / 2,
    );
    const path = [
      ...parents.map(
        (point) =>
          `M ${point.x + CARD_WIDTH / 2} ${point.y + CARD_HEIGHT} V ${lane}`,
      ),
      `M ${Math.min(...xs, center)} ${lane} H ${Math.max(...xs, center)}`,
      ...children.map(
        (point) => `M ${point.x + CARD_WIDTH / 2} ${lane} V ${point.y}`,
      ),
    ].join(" ");
    edges.push({
      id: `family:${family.id}`,
      source: family.parents[0],
      target: family.children[0],
      type: "family",
      data: { path },
      selectable: false,
      style: { stroke: "#8a969e", strokeWidth: 1.6 },
    });
  });
  return edges;
}
