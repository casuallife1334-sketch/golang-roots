import { Handle, Position, type NodeProps } from "@xyflow/react";
import {
  fullName,
  initials,
  personNodeSubtitle,
  type PersonGraphNode,
} from "../graph";

export function PersonNode({ data }: NodeProps<PersonGraphNode>) {
  const { person, photoUrl, selected } = data;
  return (
    <div className={`xy-person-node${selected ? " selected" : ""}`}>
      <Handle type="target" position={Position.Top} />
      <Handle type="source" position={Position.Bottom} />
      <Handle type="source" position={Position.Right} id="spouse-out" />
      <Handle type="target" position={Position.Left} id="spouse-in" />
      {photoUrl ? (
        <img className="xy-person-photo" src={photoUrl} alt="" />
      ) : (
        <span className="xy-person-photo xy-person-initials">
          {initials(person)}
        </span>
      )}
      <div className="xy-person-copy">
        <strong>{fullName(person)}</strong>
        <span>{personNodeSubtitle(person) || "Без дат"}</span>
        {!person.death_date && <em>Жив(а)</em>}
      </div>
    </div>
  );
}
