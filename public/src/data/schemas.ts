import { z } from "zod";

export const userSchema = z.object({
  id: z.string().min(1),
  email: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});
export const treeSchema = z.object({
  id: z.string().min(1),
  owner_id: z.string(),
  role: z.enum(["owner", "editor", "viewer"]).optional(),
  name: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});
export const personSchema = z.object({
  id: z.string().min(1),
  first_name: z.string(),
  patronymic: z.string().nullish(),
  last_name: z.string(),
  birth_date: z.string().nullish(),
  death_date: z.string().nullish(),
  gender: z.enum(["male", "female", "other"]).nullish(),
  photo_url: z.string().nullish(),
  metadata: z.record(z.string().nullable()).nullable().default(null),
  created_at: z.string(),
  updated_at: z.string().nullish(),
});
export const relationshipSchema = z.object({
  id: z.string().min(1),
  person1_id: z.string(),
  person2_id: z.string(),
  type: z.enum(["parent_child", "spouse"]),
  direction: z.enum(["parent", "child"]).nullish(),
  metadata: z.record(z.unknown()).default({}),
  created_at: z.string(),
  updated_at: z.string().nullish(),
});
export const documentSchema = z.object({
  id: z.string().min(1),
  tree_id: z.string(),
  owner: z.object({
    type: z.enum(["person", "relationship"]),
    id: z.string(),
  }),
  file_name: z.string(),
  content_type: z.string(),
  size_bytes: z.number(),
  created_by: z.string(),
  created_at: z.string(),
});
// Go list endpoints may encode an empty slice as null.
export const listOf = <T extends z.ZodTypeAny>(schema: T) =>
  z
    .array(schema)
    .nullable()
    .transform((value) => value ?? []);
export const loginSchema = z.object({
  access_token: z.string().min(1),
  token_type: z.string(),
  expires_in: z.number(),
});
