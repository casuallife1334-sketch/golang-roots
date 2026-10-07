import { useMemo, useState } from "react";
import { Grid2X2, List, Plus, Search, UsersRound, X } from "lucide-react";
import { useNavigate, useParams } from "react-router-dom";
import { useAuth } from "../auth";
import { useFamiliesData } from "../data/queries";
import { Button, EmptyState, Select, Notice } from "../shared/ui";
import { PersonPortrait } from "../shared/PersonPortrait";
import type { Family, Person } from "../types";
import { fullName } from "../utils";
import { FamilyDialog } from "../features/families/FamilyDialog";

type FamilyCard = {
  id: string;
  members: Person[];
  adults: Person[];
  children: Person[];
  name: string;
};

export function FamiliesPage() {
  const { treeId } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const { trees, tree, people, families } = useFamiliesData(treeId);
  const [openedFamily, setOpenedFamily] = useState<string | "new" | null>(null);
  const [search, setSearch] = useState("");
  const [mode, setMode] = useState<"grid" | "list">("grid");
  const [filter, setFilter] = useState("all");
  const [sort, setSort] = useState("name");
  const familyCards = useMemo(
    () => (families.data ?? []).map((family) => toFamilyCard(family, people.data ?? [])),
    [families.data, people.data],
  );
  const filteredFamilies = useMemo(
    () =>
      familyCards
        .filter(
          (family) =>
            filter === "all" ||
            (filter === "children" && family.children.length > 0) ||
            (filter === "empty" && family.children.length === 0),
        )
        .filter((family) =>
          `${family.name} ${family.members.map(fullName).join(" ")}`
            .toLowerCase()
            .includes(search.toLowerCase()),
        )
        .sort((a, b) =>
          sort === "children"
            ? b.children.length - a.children.length
            : a.name.localeCompare(b.name),
        ),
    [familyCards, search, filter, sort],
  );

  if (
    trees.isPending ||
    (tree && (people.isPending || families.isPending))
  )
    return <div className="center-panel">Загрузка семей…</div>;
  if (trees.error || people.error || families.error)
    return (
      <div className="center-panel">
        <Notice>
          {trees.error?.message ||
            people.error?.message ||
            families.error?.message}
        </Notice>
        <Button
          onClick={() => {
            void trees.refetch();
            void people.refetch();
            void families.refetch();
          }}
        >
          Повторить
        </Button>
      </div>
    );
  if (!tree)
    return (
      <EmptyState
        title="Дерево не найдено"
        text="Оно удалено или недоступно"
        action={<Button onClick={() => navigate("/trees")}>К деревьям</Button>}
      />
    );
  const canWrite = tree.role
    ? tree.role === "owner" || tree.role === "editor"
    : tree.owner_id === user?.id;
  const selectedFamily = families.data?.find((family) => family.id === openedFamily);
  return (
    <div className="content-page families-page">
      <header className="topbar">
        <div className="search-box">
          <Search size={18} />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Поиск по семьям и участникам…"
            aria-label="Поиск по семьям и участникам"
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
            {filteredFamilies.length} {filteredFamilies.length === 1 ? "семья" : "семей"} в
            выбранном дереве
          </p>
        </div>
        <div className="family-heading-actions">
          {canWrite && (
            <Button onClick={() => setOpenedFamily("new")}>
              <Plus size={16} />
              Создать семью
            </Button>
          )}
          <Button variant="secondary" onClick={() => navigate(`/trees/${tree.id}`)}>
            <UsersRound size={16} />
            К дереву
          </Button>
        </div>
      </section>
      {filteredFamilies.length ? (
        <div className={`family-grid ${mode}`}>
          {filteredFamilies.map((family) => {
            const portraits = family.adults.length ? family.adults : family.members;
            return (
              <button
                className="family-card"
                key={family.id}
                onClick={() => setOpenedFamily(family.id)}
                aria-label={`Открыть семью: ${family.name}`}
              >
                <div className="family-photos">
                  {portraits.slice(0, 2).map((person) => (
                    <PersonPortrait
                      person={person}
                      treeId={treeId!}
                      className="portrait"
                      key={person.id}
                    />
                  ))}
                </div>
                <div className="family-card-info">
                  <h3>{family.name}</h3>
                  <small>{childrenLabel(family.children.length)}</small>
                </div>
              </button>
            );
          })}
        </div>
      ) : (
        <EmptyState
          icon={<UsersRound size={28} />}
          title={familyCards.length ? "Семьи не найдены" : "Семей пока нет"}
          text={familyCards.length ? "Попробуйте изменить поиск или фильтр." : "Создайте семью или добавьте связь между людьми в выбранном дереве."}
          action={familyCards.length
            ? <Button onClick={() => { setSearch(""); setFilter("all"); }}>Сбросить фильтры</Button>
            : canWrite
              ? <Button onClick={() => setOpenedFamily("new")}>Создать семью</Button>
              : <Button onClick={() => navigate(`/trees/${tree.id}`)}>Открыть древо</Button>}
        />
      )}
      {(openedFamily === "new" || selectedFamily) && (
        <FamilyDialog
          key={openedFamily}
          tree={tree}
          family={selectedFamily ?? null}
          people={people.data ?? []}
          editable={canWrite}
          onClose={() => setOpenedFamily(null)}
          onCreated={(id) => setOpenedFamily(id)}
          onDeleted={() => setOpenedFamily(null)}
          onOpenPerson={(id) => navigate(`/trees/${tree.id}?person=${id}`)}
        />
      )}
    </div>
  );
}

function toFamilyCard(family: Family, people: Person[]): FamilyCard {
  const peopleByID = new Map(people.map((person) => [person.id, person]));
  const memberIDs = [...new Set(family.members.map((member) => member.person_id))];
  const members = memberIDs
    .map((id) => peopleByID.get(id))
    .filter((person): person is Person => Boolean(person));
  const adultIDs = new Set(
    family.members
      .filter((member) => member.role !== "child")
      .map((member) => member.person_id),
  );
  const childrenIDs = new Set(
    family.members
      .filter((member) => member.role === "child")
      .map((member) => member.person_id),
  );
  const adults = members.filter((person) => adultIDs.has(person.id));
  const children = members.filter((person) => childrenIDs.has(person.id) && !adultIDs.has(person.id));
  return {
    id: family.id,
    members,
    adults,
    children,
    name: family.name.trim() || familyName(adults.length ? adults : members),
  };
}

function familyName(members: Person[]) {
  return members.length ? members.map(fullName).join(" и ") : "Новая семья";
}

function childrenLabel(count: number) {
  const lastTwo = count % 100;
  const last = count % 10;
  if (lastTwo >= 11 && lastTwo <= 14) return `${count} детей`;
  if (last === 1) return `${count} ребенок`;
  if (last >= 2 && last <= 4) return `${count} ребенка`;
  return `${count} детей`;
}
