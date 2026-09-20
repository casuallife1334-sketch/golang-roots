import { useEffect, useState } from "react";
import { useQueries } from "@tanstack/react-query";
import { api } from "../api";
import type { Person } from "../types";

export function usePersonPhotoUrls(treeId: string, people: Person[]) {
  const queries = useQueries({
    queries: people.map((person) => ({
      queryKey: ["family-photo", treeId, person.id, person.photo_url],
      queryFn: () => api.photo(treeId, person.id),
      enabled: Boolean(treeId && person.photo_url),
      staleTime: Infinity,
    })),
  });
  const [urls, setUrls] = useState<Record<string, string>>({});
  const ids = people.map((person) => person.id).join("|");
  const blobs = queries.map((query) => query.data);
  useEffect(() => {
    const next: Record<string, string> = {};
    people.forEach((person, index) => {
      const blob = blobs[index];
      if (blob) next[person.id] = URL.createObjectURL(blob);
    });
    setUrls(next);
    return () => Object.values(next).forEach((url) => URL.revokeObjectURL(url));
  }, [treeId, ids, ...blobs]);
  return urls;
}
