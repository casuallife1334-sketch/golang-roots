import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Grid2X2, List, Search, UsersRound, X } from "lucide-react";
import { useNavigate, useParams } from "react-router-dom";
import { api } from "../api";
import { useAuth } from "../auth";
import { Button, EmptyState, Select } from "../shared/ui";
import { PersonPortrait } from "../shared/PersonPortrait";
import type { Person, Relationship } from "../types";
import { formatYears, fullName, initials } from "../utils";

interface Family {
  id: string;
  members: Person[];
  children: Person[];
  name: string;
}
function deriveFamilies(
  people: Person[],
  relationships: Relationship[],
): Family[] {
  const map = new Map(people.map((p) => [p.id, p]));
  const groups = new Map<string, Family>();
  for (const r of relationships.filter((r) => r.type === "spouse")) {
    const ids = [r.person1_id, r.person2_id].sort();
    const members = ids.map((id) => map.get(id)).filter(Boolean) as Person[];
    if (members.length === 2)
      groups.set(ids.join(":"), {
        id: ids.join(":"),
        members,
        children: [],
        name: familyName(members),
      });
  }
  const parentSets = new Map<string, string[]>();
  relationships
    .filter((r) => r.type === "parent_child")
    .forEach((r) => {
      const parent = r.direction === "child" ? r.person2_id : r.person1_id;
      const child = r.direction === "child" ? r.person1_id : r.person2_id;
      parentSets.set(child, [...(parentSets.get(child) || []), parent]);
    });
  parentSets.forEach((parentIds, childId) => {
    const ids = [...new Set(parentIds)].sort();
    const members = ids.map((id) => map.get(id)).filter(Boolean) as Person[];
    if (!members.length) return;
    const key = ids.join(":");
    const current = groups.get(key) || {
      id: key,
      members,
      children: [],
      name:
        members.length > 1 ? familyName(members) : "Один известный родитель",
    };
    current.children.push(map.get(childId)!);
    groups.set(key, current);
  });
  return [...groups.values()]
    .filter((f) => f.members.length)
    .sort((a, b) => a.name.localeCompare(b.name));
}
function familyName(members: Person[]) {
  const surnames = [
    ...new Set(members.map((p) => p.last_name).filter(Boolean)),
  ];
  return surnames.length === 1
    ? `Семья ${surnames[0]}`
    : members.map((p) => p.first_name).join(" и ");
}
export function FamiliesPage() {
  const { treeId } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const trees = useQuery({ queryKey: ["trees", user?.id], queryFn: api.trees });
  const tree = trees.data?.find((t) => t.id === treeId);
  const people = useQuery({
    queryKey: ["people", user?.id, treeId],
    queryFn: () => api.persons(treeId!),
    enabled: Boolean(treeId),
  });
  const relationships = useQuery({
    queryKey: ["relationships", user?.id, treeId],
    queryFn: () => api.relationships(treeId!),
    enabled: Boolean(treeId),
  });
  const [search, setSearch] = useState("");
  const [mode, setMode] = useState<"grid" | "list">("grid");
  const [filter, setFilter] = useState("all");
  const [sort, setSort] = useState("name");
  const families = useMemo(
    () =>
      deriveFamilies(people.data || [], relationships.data || [])
        .filter(
          (f) =>
            filter === "all" ||
            (filter === "children" && f.children.length) ||
            (filter === "empty" && !f.children.length),
        )
        .filter((f) =>
          `${f.name} ${f.members.map(fullName).join(" ")}`
            .toLowerCase()
            .includes(search.toLowerCase()),
        )
        .sort((a, b) =>
          sort === "children"
            ? b.children.length - a.children.length
            : a.name.localeCompare(b.name),
        ),
    [people.data, relationships.data, search, filter, sort],
  );
  if (!tree) return <div className="center-panel">Загрузка семей…</div>;
  return (
    <div className="content-page">
      <header className="topbar">
        <div className="search-box">
          <Search size={18} />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Поиск по семьям и участникам…"
          />
          {search && (
            <button onClick={() => setSearch("")} aria-label="Очистить">
              <X size={15} />
            </button>
          )}
        </div>
        <div className="family-controls">
          <Select value={sort} onChange={(e) => setSort(e.target.value)}>
            <option value="name">По имени</option>
            <option value="children">По количеству детей</option>
          </Select>
          <Select value={filter} onChange={(e) => setFilter(e.target.value)}>
            <option value="all">Все семьи</option>
            <option value="children">С детьми</option>
            <option value="empty">Без детей</option>
          </Select>
          <button
            className={`icon-button ${mode === "grid" ? "active-control" : ""}`}
            onClick={() => setMode("grid")}
            aria-label="Сетка"
          >
            <Grid2X2 size={17} />
          </button>
          <button
            className={`icon-button ${mode === "list" ? "active-control" : ""}`}
            onClick={() => setMode("list")}
            aria-label="Список"
          >
            <List size={17} />
          </button>
        </div>
      </header>
      <section className="page-heading">
        <div>
          <h1>Семьи</h1>
          <p>
            {families.length} {families.length === 1 ? "семья" : "семей"} в
            выбранном дереве
          </p>
        </div>
        <Button onClick={() => navigate(`/trees/${treeId}`)}>
          <UsersRound size={16} />
          Управлять связями
        </Button>
      </section>
      {people.isLoading ? (
        <div className="skeleton-grid">
          {[1, 2, 3].map((x) => (
            <div className="skeleton-card" key={x} />
          ))}
        </div>
      ) : families.length ? (
        <div className={`family-grid ${mode}`}>
          {families.map((family) => (
            <button
              className="family-card"
              key={family.id}
              onClick={() =>
                navigate(`/trees/${treeId}?person=${family.members[0].id}`)
              }
            >
              <div className="family-photos">
                {family.members.slice(0, 2).map((p) => (
                  <PersonPortrait
                    person={p}
                    treeId={treeId!}
                    className="portrait"
                    key={p.id}
                  />
                ))}
              </div>
              <h3>{family.name}</h3>
              <p>{family.members.map(fullName).join(" · ")}</p>
              <small>
                {family.children.length}{" "}
                {family.children.length === 1 ? "ребёнок" : "детей"}
              </small>
            </button>
          ))}
        </div>
      ) : (
        <EmptyState
          icon={<UsersRound size={28} />}
          title="Семьи появятся после добавления связей"
          text="Создайте связь между двумя людьми в выбранном дереве."
          action={
            <Button onClick={() => navigate(`/trees/${treeId}`)}>
              Открыть древо
            </Button>
          }
        />
      )}
    </div>
  );
}
