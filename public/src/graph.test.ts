import { describe, expect, it } from "vitest";
import { buildGraph } from "./graph";
import type { Person, Relationship } from "./types";

const person = (id: string): Person => ({
  id,
  first_name: id,
  last_name: "Person",
  metadata: null,
  created_at: "2026-01-01T00:00:00Z",
});
const relation = (
  id: string,
  person1_id: string,
  person2_id: string,
  type: Relationship["type"] = "spouse",
): Relationship => ({ id, person1_id, person2_id, type, created_at: "" });

describe("buildGraph", () => {
  it("keeps people without relationships visible", () =>
    expect(
      buildGraph([person("root"), person("other")], [], {}, "root").nodes,
    ).toHaveLength(2));
  it("creates a junction for parent links", () => {
    const result = buildGraph(
      [person("parent"), person("child")],
      [relation("parent-child", "parent", "child", "parent_child")],
      {},
      "parent",
    );
    expect(result.nodes.some((node) => node.type === "junction")).toBe(true);
    expect(result.edges).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ source: "parent", target: "junction-0" }),
        expect.objectContaining({ source: "junction-0", target: "child" }),
      ]),
    );
  });
  it("keeps coordinates stable when selection changes", () => {
    const people = [person("parent"), person("child")];
    const relationships = [
      relation("parent-child", "parent", "child", "parent_child"),
    ];
    const first = buildGraph(people, relationships, {}, "parent");
    const second = buildGraph(people, relationships, {}, "child");
    expect(second.nodes.map((node) => [node.id, node.position])).toEqual(
      first.nodes.map((node) => [node.id, node.position]),
    );
  });
  it("keeps parents of spouses on separate family junctions", () => {
    const people = [
      person("mother-a"),
      person("father-a"),
      person("mother-b"),
      person("father-b"),
      person("child-a"),
      person("child-b"),
    ];
    const relationships = [
      relation("spouse-a", "mother-a", "father-a"),
      relation("spouse-b", "mother-b", "father-b"),
      relation("spouse-children", "child-a", "child-b"),
      relation("a-child", "mother-a", "child-a", "parent_child"),
      relation("a-child-2", "father-a", "child-a", "parent_child"),
      relation("b-child", "mother-b", "child-b", "parent_child"),
      relation("b-child-2", "father-b", "child-b", "parent_child"),
    ];
    const result = buildGraph(people, relationships, {}, "child-a");
    expect(
      result.nodes.filter((node) => node.type === "junction"),
    ).toHaveLength(2);
  });
});
