import { renderToStaticMarkup } from "react-dom/server";
import { expect, it, vi } from "vitest";
import { FamiliesPage } from "./FamiliesPage";
import { useFamiliesData } from "../data/queries";
import type { Family, Person } from "../types";

vi.mock("../data/queries", () => ({ useFamiliesData: vi.fn() }));
vi.mock("../auth", () => ({ useAuth: () => ({ user: { id: "user" } }) }));
vi.mock("react-router-dom", async (importOriginal) => ({
  ...(await importOriginal<typeof import("react-router-dom")>()),
  useParams: () => ({ treeId: "tree" }),
  useNavigate: () => vi.fn(),
}));
vi.mock("../shared/PersonPortrait", () => ({
  PersonPortrait: () => <span>Портрет</span>,
}));

const person = {
  id: "parent",
  first_name: "Анна",
  last_name: "Попович",
} as Person;
const family = {
  id: "family",
  members: [{ person_id: person.id, role: "parent" }],
  name: "",
} as Family;

function setData(families: Family[]) {
  vi.mocked(useFamiliesData).mockReturnValue({
    trees: { isPending: false, data: [{ id: "tree" }] },
    tree: { id: "tree" },
    people: { isPending: false, data: [person] },
    families: { isPending: false, data: families },
  } as unknown as ReturnType<typeof useFamiliesData>);
}

it("shows an empty state when the families API returns an empty list", () => {
  setData([]);
  const html = renderToStaticMarkup(<FamiliesPage />);
  expect(html).toContain("Семей пока нет");
  expect(html).not.toContain("Анна Попович");
});

it("shows families returned by the API", () => {
  setData([family]);
  const html = renderToStaticMarkup(<FamiliesPage />);
  expect(html).toContain("Анна Попович");
  expect(html).not.toContain("Семей пока нет");
});
