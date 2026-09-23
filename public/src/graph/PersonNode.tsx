import { Handle, Position, type NodeProps } from "@xyflow/react";
import type { PersonGraphNode } from "../graph";
import { fullName, initials, formatYears } from "../utils";
import { PersonPortrait } from "../shared/PersonPortrait";
export function PersonNode({ data }: NodeProps<PersonGraphNode>) {
  const { person, treeId, selected, portraits } = data;
  return (
    <div
      className={`xy-person-node${selected ? " selected" : ""}`}
      title={fullName(person)}
    >
      <Handle type="target" position={Position.Top} />
      <Handle type="source" position={Position.Bottom} />
      <Handle type="source" position={Position.Right} id="spouse-out" />
      <Handle type="target" position={Position.Left} id="spouse-in" />
      {portraits ? (
        <PersonPortrait
          person={person}
          treeId={treeId}
          className="xy-person-photo"
        />
      ) : (
        <span className="xy-person-photo xy-person-initials">
          {initials(person)}
        </span>
      )}
      <div className="xy-person-copy">
        <strong>{fullName(person)}</strong>
        <span>{formatYears(person)}</span>
        {person.metadata?.city && <span>{person.metadata.city}</span>}
        {!person.death_date && <em>Жив(а)</em>}
      </div>
    </div>
  );
}
