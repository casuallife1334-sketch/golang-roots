import { describe, expect, it } from "vitest";
import { buildGraph, CARD_HEIGHT, CARD_WIDTH } from "./graph";
import { familyModel, relationshipProblem } from "./graph/model";
import type { Person, Relationship } from "./types";
const person = (id: string): Person => ({
  id,
  first_name: id,
  last_name: "Person",
  metadata: null,
  created_at: "",
});
const rel = (
  id: string,
  a: string,
  b: string,
  type: Relationship["type"] = "parent_child",
): Relationship => ({
  id,
  person1_id: a,
  person2_id: b,
  type,
  metadata: {},
  created_at: "",
});
const people = ["a", "b", "c", "d", "x", "y", "z", "isolated"].map(person);
const relations = [
  rel("ab", "a", "b", "spouse"),
  rel("cd", "c", "d", "spouse"),
  rel("xy", "x", "y", "spouse"),
  rel("ax", "a", "x"),
  rel("bx", "b", "x"),
  rel("cy", "c", "y"),
  rel("dy", "d", "y"),
  rel("xz", "x", "z"),
  rel("yz", "y", "z"),
];
describe("family graph", () => {
  it("does not combine parents of married children", () => {
    expect(familyModel(people, relations).families).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ parents: ["a", "b"], children: ["x"] }),
        expect.objectContaining({ parents: ["c", "d"], children: ["y"] }),
      ]),
    );
    expect(
      buildGraph(people, relations).edges.filter(
        (edge) => edge.type === "family",
      ),
    ).toHaveLength(3);
  });
  it("is independent of API order and preserves isolated people", () => {
    expect(buildGraph([...people].reverse(), [...relations].reverse())).toEqual(
      buildGraph(people, relations),
    );
    expect(
      buildGraph(people, relations).nodes.map((node) => node.id),
    ).toContain("isolated");
  });
  it("places parents above children without overlapping cards", () => {
    const { nodes } = buildGraph(people, relations);
    for (const a of nodes)
      for (const b of nodes) {
        if (a.id === b.id) continue;
        expect(
          Math.abs(a.position.x - b.position.x) >= CARD_WIDTH ||
            Math.abs(a.position.y - b.position.y) >= CARD_HEIGHT,
        ).toBe(true);
      }
    for (const relation of relations.filter(
      (item) => item.type === "parent_child",
    )) {
      expect(
        nodes.find((node) => node.id === relation.person2_id)!.position.y,
      ).toBeGreaterThan(
        nodes.find((node) => node.id === relation.person1_id)!.position.y,
      );
    }
  });
  it("does not transitively merge multiple marriages", () => {
    expect(
      familyModel(people, [
        rel("ab", "a", "b", "spouse"),
        rel("ac", "a", "c", "spouse"),
      ]).spouses,
    ).toHaveLength(2);
    expect(
      buildGraph(people, [
        rel("ab", "a", "b", "spouse"),
        rel("ac", "a", "c", "spouse"),
      ]).nodes,
    ).toHaveLength(people.length);
  });
  it("orders married children consistently with their separate parent branches", () => {
    const graph = buildGraph(people, relations);
    const x = (id: string) =>
      graph.nodes.find((node) => node.id === id)!.position.x;
    expect(Math.sign(x("x") - x("y"))).toBe(
      Math.sign((x("a") + x("b")) / 2 - (x("c") + x("d")) / 2),
    );
  });
  it("rejects parent cycles and ignores dangling links", () => {
    expect(
      relationshipProblem(people, relations, {
        person1_id: "z",
        person2_id: "a",
        type: "parent_child",
      }),
    ).toContain("цикл");
    const graph = buildGraph(people, [
      ...relations,
      rel("missing", "missing", "x"),
      rel("self", "z", "z"),
    ]);
    expect(graph.warnings).toHaveLength(2);
    expect(
      graph.edges.every(
        (edge) =>
          graph.nodes.some((node) => node.id === edge.source) &&
          graph.nodes.some((node) => node.id === edge.target),
      ),
    ).toBe(true);
  });
});
