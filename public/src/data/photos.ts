import { useQueries, useQuery } from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import { api } from "../api";
import { useAuth } from "../auth";
import { keys } from "./queries";
import type { Person } from "../types";

const photoOptions = (
  user: string | undefined,
  tree: string,
  person: Person,
) => ({
  queryKey: [...keys.photo(user, tree, person.id), person.photo_url],
  queryFn: ({ signal }: { signal: AbortSignal }) =>
    api.photo(tree, person.id, signal),
  enabled: Boolean(user && tree && person.photo_url),
  staleTime: 300_000,
});
function usePhotoUrls(entries: { id: string; blob?: Blob }[]) {
  const cache = useRef(new Map<string, { blob: Blob; url: string }>());
  const [urls, setUrls] = useState<Record<string, string>>({});
  useEffect(() => {
    const next: Record<string, string> = {};
    let changed = false;
    const active = new Set(
      entries.filter((entry) => entry.blob).map((entry) => entry.id),
    );
    for (const [id, value] of cache.current)
      if (!active.has(id)) {
        URL.revokeObjectURL(value.url);
        cache.current.delete(id);
        changed = true;
      }
    for (const { id, blob } of entries) {
      if (!blob) continue;
      let value = cache.current.get(id);
      if (value?.blob !== blob) {
        if (value) URL.revokeObjectURL(value.url);
        value = { blob, url: URL.createObjectURL(blob) };
        cache.current.set(id, value);
        changed = true;
      }
      next[id] = value!.url;
    }
    if (changed) setUrls(next);
  });
  useEffect(
    () => () => {
      cache.current.forEach((value) => URL.revokeObjectURL(value.url));
      cache.current.clear();
    },
    [],
  );
  return urls;
}
export function usePersonPhoto(tree: string, person: Person) {
  const { user } = useAuth();
  const query = useQuery(photoOptions(user?.id, tree, person));
  const blob = person.photo_url ? query.data : undefined;
  const urls = usePhotoUrls([{ id: person.id, blob }]);
  return { url: blob ? urls[person.id] : undefined, error: query.error };
}
export function usePersonPhotoUrls(tree: string, people: Person[]) {
  const { user } = useAuth();
  const queries = useQueries({
    queries: people.map((person) => photoOptions(user?.id, tree, person)),
  });
  return usePhotoUrls(
    people.map((person, index) => ({
      id: person.id,
      blob: person.photo_url ? queries[index].data : undefined,
    })),
  );
}
