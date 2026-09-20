import { Handle, Position, type NodeProps } from "@xyflow/react";
import type { JunctionGraphNode } from "../graph";

export function JunctionNode({}: NodeProps<JunctionGraphNode>) {
  return (
    <div className="xy-family-junction" aria-hidden="true">
      <Handle type="target" position={Position.Top} />
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
}
