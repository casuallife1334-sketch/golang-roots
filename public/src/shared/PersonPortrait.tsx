import { usePersonPhoto } from "../data/photos";
import type { Person } from "../types";
import { initials } from "../utils";

export function PersonPortrait({
  person,
  treeId,
  className = "",
}: {
  person: Person;
  treeId: string;
  className?: string;
}) {
  const { url } = usePersonPhoto(treeId, person);
  return (
    <span className={`person-portrait ${className}`}>
      {url ? <img src={url} alt="" /> : initials(person)}
    </span>
  );
}
