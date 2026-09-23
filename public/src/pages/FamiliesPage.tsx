import { useMemo, useState } from "react";
import { Grid2X2, List, Search, UsersRound, X } from "lucide-react";
import { useNavigate, useParams } from "react-router-dom";
import { useTreeData } from "../data/queries";
import { Button, EmptyState, Select, Notice } from "../shared/ui";
import { PersonPortrait } from "../shared/PersonPortrait";
import { fullName } from "../utils";
import { deriveFamilies } from "../features/tree/families";

export function FamiliesPage() {
  const { treeId } = useParams();
  const navigate = useNavigate();
  const { trees, tree, people, relationships } = useTreeData(treeId);
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
  if (
    trees.isPending ||
    (tree && (people.isPending || relationships.isPending))
  )
    return <div className="center-panel">Загрузка семей…</div>;
  if (trees.error || people.error || relationships.error)
    return (
      <div className="center-panel">
        <Notice>
          {trees.error?.message ||
            people.error?.message ||
            relationships.error?.message}
        </Notice>
        <Button
          onClick={() => {
            void trees.refetch();
            void people.refetch();
            void relationships.refetch();
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
