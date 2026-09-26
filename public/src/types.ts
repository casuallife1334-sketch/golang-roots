export type Gender = "male" | "female" | "other";
export type RelationshipType = "parent_child" | "spouse";
export type RelationshipDirection = "parent" | "child";
export type DocumentOwnerType = "person" | "relationship";

export interface User {
  id: string;
  email: string;
  created_at: string;
  updated_at: string;
}
export interface Tree {
  id: string;
  owner_id: string;
  role?: "owner" | "editor" | "viewer";
  name: string;
  created_at: string;
  updated_at: string;
}
export interface Person {
  id: string;
  first_name: string;
  patronymic?: string | null;
  last_name: string;
  birth_date?: string | null;
  death_date?: string | null;
  gender?: Gender | null;
  photo_url?: string | null;
  metadata: Record<string, string | null> | null;
  created_at: string;
  updated_at?: string | null;
}
export interface Relationship {
  id: string;
  person1_id: string;
  person2_id: string;
  type: RelationshipType;
  direction?: RelationshipDirection | null;
  metadata: Record<string, unknown>;
  created_at: string;
  updated_at?: string | null;
}
export interface PersonInput {
  first_name: string;
  patronymic?: string | null;
  last_name: string;
  birth_date?: string | null;
  death_date?: string | null;
  gender?: Gender | null;
  metadata: Record<string, unknown>;
}
export interface RelationshipInput {
  person1_id: string;
  person2_id: string;
  type: RelationshipType;
  direction?: RelationshipDirection;
  metadata?: Record<string, unknown>;
}
export interface PatchRelationshipInput {
  metadata?: Record<string, unknown>;
}
export interface DocumentOwner {
  type: DocumentOwnerType;
  id: string;
}
export interface Document {
  id: string;
  tree_id: string;
  owner: DocumentOwner;
  file_name: string;
  content_type: string;
  size_bytes: number;
  created_by: string;
  created_at: string;
}

export const EMPTY_ARRAY = <T>(value: T[] | null | undefined): T[] =>
  value ?? [];
