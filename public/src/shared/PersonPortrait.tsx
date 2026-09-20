import { useQuery } from "@tanstack/react-query";
import { api } from "../api";
import type { Person } from "../types";
import { initials, useObjectUrl } from "../utils";

export function PersonPortrait({
  person,
  treeId,
  className = "",
}: {
  person: Person;
  treeId: string;
  className?: string;
}) {
  const photo = useQuery({
    queryKey: ["photo", treeId, person.id, person.photo_url],
    queryFn: () => api.photo(treeId, person.id),
    enabled: Boolean(person.photo_url),
    staleTime: Infinity,
  });
  const url = useObjectUrl(photo.data);
  return (
    <span className={`person-portrait ${className}`}>
      {url ? <img src={url} alt="" /> : initials(person)}
    </span>
  );
}
