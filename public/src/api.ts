import type { PersonInput, RelationshipInput } from "./types";
import { request } from "./data/http";
import {
  userSchema,
  treeSchema,
  personSchema,
  relationshipSchema,
  loginSchema,
  listOf,
} from "./data/schemas";
export { ApiError, setUnauthorizedHandler } from "./data/http";
const json = (method: string, body: unknown): RequestInit => ({
  method,
  body: JSON.stringify(body),
});

export const api = {
  register: (body: { email: string; password: string }) =>
    request("/auth/register", userSchema, json("POST", body), true),
  login: (body: { email: string; password: string }) =>
    request("/auth/login", loginSchema, json("POST", body), true),
  me: (signal?: AbortSignal) => request("/users/me", userSchema, { signal }),
  trees: (signal?: AbortSignal) =>
    request("/trees", listOf(treeSchema), { signal }),
  tree: (id: string, signal?: AbortSignal) =>
    request(`/trees/${id}`, treeSchema, { signal }),
  createTree: (name: string) =>
    request("/trees", treeSchema, json("POST", { name })),
  patchTree: (id: string, name: string) =>
    request(`/trees/${id}`, treeSchema, json("PATCH", { name })),
  deleteTree: (id: string) =>
    request<void>(`/trees/${id}`, "void", { method: "DELETE" }),
  persons: (treeId: string, signal?: AbortSignal) =>
    request(`/trees/${treeId}/persons`, listOf(personSchema), { signal }),
  person: (treeId: string, id: string) =>
    request(`/trees/${treeId}/persons/${id}`, personSchema),
  createPerson: (treeId: string, body: PersonInput) =>
    request(`/trees/${treeId}/persons`, personSchema, json("POST", body)),
  patchPerson: (treeId: string, id: string, body: Partial<PersonInput>) =>
    request(
      `/trees/${treeId}/persons/${id}`,
      personSchema,
      json("PATCH", body),
    ),
  deletePerson: (treeId: string, id: string) =>
    request<void>(`/trees/${treeId}/persons/${id}`, "void", {
      method: "DELETE",
    }),
  relationships: (treeId: string, signal?: AbortSignal) =>
    request(`/trees/${treeId}/relationships`, listOf(relationshipSchema), {
      signal,
    }),
  createRelationship: (treeId: string, body: RelationshipInput) =>
    request(
      `/trees/${treeId}/relationships`,
      relationshipSchema,
      json("POST", body),
    ),
  deleteRelationship: (treeId: string, id: string) =>
    request<void>(`/trees/${treeId}/relationships/${id}`, "void", {
      method: "DELETE",
    }),
  photo: (treeId: string, id: string, signal?: AbortSignal) =>
    request<Blob>(`/trees/${treeId}/persons/${id}/photo`, "blob", { signal }),
  uploadPhoto: (treeId: string, id: string, file: File) => {
    const body = new FormData();
    body.append("file", file);
    return request(`/trees/${treeId}/persons/${id}/photo`, personSchema, {
      method: "POST",
      body,
    });
  },
  deletePhoto: (treeId: string, id: string) =>
    request<void>(`/trees/${treeId}/persons/${id}/photo`, "void", {
      method: "DELETE",
    }),
};
