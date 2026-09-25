import { queryOptions, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "../api";
import { useAuth } from "../auth";
import type { Person, Relationship } from "../types";

export const keys = {
  trees: (user?: string) => ["trees", user] as const,
  people: (user?: string, tree?: string) => ["people", user, tree] as const,
  relationships: (user?: string, tree?: string) =>
    ["relationships", user, tree] as const,
  photo: (user?: string, tree?: string, person?: string) =>
    ["photo", user, tree, person] as const,
};
export const treeQuery = (user?: string) =>
  queryOptions({
    queryKey: keys.trees(user),
    queryFn: ({ signal }) => api.trees(signal),
    enabled: Boolean(user),
  });
export const peopleQuery = (user?: string, tree?: string) =>
  queryOptions({
    queryKey: keys.people(user, tree),
    queryFn: ({ signal }) => api.persons(tree!, signal),
    enabled: Boolean(user && tree),
  });
export const relationshipsQuery = (user?: string, tree?: string) =>
  queryOptions({
    queryKey: keys.relationships(user, tree),
    queryFn: ({ signal }) => api.relationships(tree!, signal),
    enabled: Boolean(user && tree),
  });
export function useTreeData(treeId?: string) {
  const { user } = useAuth();
  const trees = useQuery(treeQuery(user?.id));
  const tree = treeId
    ? trees.data?.find((item) => item.id === treeId)
    : trees.data?.[0];
  const people = useQuery(peopleQuery(user?.id, tree?.id));
  const relationships = useQuery(relationshipsQuery(user?.id, tree?.id));
  return { trees, tree, people, relationships };
}
export function useTreeCache(treeId: string) {
  const { user } = useAuth();
  const query = useQueryClient();
  return {
    savePerson: (person: Person) => {
      query.setQueryData<Person[]>(keys.people(user?.id, treeId), (old) => [
        ...(old ?? []).filter((item) => item.id !== person.id),
        person,
      ]);
    },
    saveRelationship: (relationship: Relationship) => {
      query.setQueryData<Relationship[]>(
        keys.relationships(user?.id, treeId),
        (old) => [
          ...(old ?? []).filter((item) => item.id !== relationship.id),
          relationship,
        ],
      );
    },
    refresh: () =>
      Promise.all([
        query.invalidateQueries({ queryKey: keys.people(user?.id, treeId) }),
        query.invalidateQueries({
          queryKey: keys.relationships(user?.id, treeId),
        }),
      ]),
    refreshPhoto: (id: string) =>
      query.invalidateQueries({ queryKey: keys.photo(user?.id, treeId, id) }),
  };
}
