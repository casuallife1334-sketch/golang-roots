import { expect, it } from "vitest";
import type { Family, Person } from "../../types";
import { familyOptionLabel } from "./familyOptionLabel";

const people = [
  { id: "alexey", first_name: "Алексей", last_name: "Попович" },
  { id: "anna", first_name: "Анна", last_name: "Попович" },
] as Person[];

it("shows each parent once even when they have both partner and parent roles", () => {
  const family = {
    id: "family",
    name: "",
    members: [
      { person_id: "alexey", role: "partner" },
      { person_id: "alexey", role: "parent" },
      { person_id: "anna", role: "partner" },
      { person_id: "anna", role: "parent" },
      { person_id: "child", role: "child" },
    ],
  } as Family;

  expect(familyOptionLabel(family, people)).toBe("Алексей Попович и Анна Попович");
});

it("uses a custom name when one is set", () => {
  const family = { id: "family", name: "Моя семья", members: [] } as unknown as Family;
  expect(familyOptionLabel(family, people)).toBe("Моя семья");
});
