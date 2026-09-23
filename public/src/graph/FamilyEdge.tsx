import { BaseEdge, type EdgeProps } from "@xyflow/react";
export function FamilyEdge({ id, data, style }: EdgeProps) {
  return (
    <BaseEdge
      id={id}
      path={String(data?.path ?? "")}
      style={{ ...style, strokeLinejoin: "round", strokeLinecap: "round" }}
    />
  );
}
