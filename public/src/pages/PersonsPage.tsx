import { useQuery } from "@tanstack/react-query";
import { ContactRound, Search, X } from "lucide-react";
import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useAuth } from "../auth";
import { peopleQuery, treeQuery } from "../data/queries";
import { Button, EmptyState, Notice } from "../shared/ui";
import { formatDate, fullName } from "../utils";

export function PersonsPage() {
  const { treeId } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [search, setSearch] = useState("");
  const trees = useQuery(treeQuery(user?.id));
  const tree = trees.data?.find((item) => item.id === treeId);
  const people = useQuery(peopleQuery(user?.id, tree?.id));

  if (trees.isPending || (tree && people.isPending)) {
    return <div className="center-panel">Загрузка персон…</div>;
  }

  if (trees.error || people.error) {
    return (
      <div className="center-panel">
        <Notice>{trees.error?.message || people.error?.message}</Notice>
        <Button
          onClick={() => {
            void trees.refetch();
            void people.refetch();
          }}
        >
          Повторить
        </Button>
      </div>
    );
  }

  if (!tree) {
    return (
      <EmptyState
        icon={<ContactRound size={28} />}
        title="Дерево не найдено"
        text="Оно удалено или недоступно"
        action={<Button onClick={() => navigate("/trees")}>К деревьям</Button>}
      />
    );
  }

  const items = people.data ?? [];
  const normalizedSearch = search.trim().toLocaleLowerCase();
  const filteredItems = normalizedSearch
    ? items.filter((person) =>
        fullName(person).toLocaleLowerCase().includes(normalizedSearch),
      )
    : items;

  return (
    <div className="content-page persons-page">
      <header className="topbar">
        <div className="search-box">
          <Search size={18} />
          <input
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Поиск по ФИО…"
            aria-label="Поиск по ФИО"
          />
          {search && (
            <button onClick={() => setSearch("")} aria-label="Очистить поиск">
              <X size={15} />
            </button>
          )}
        </div>
      </header>
      <section className="page-heading">
        <div>
          <h1>Персоны</h1>
          <p>
            {items.length} {peopleCountLabel(items.length)} в дереве «{tree.name}»
          </p>
        </div>
      </section>
      {filteredItems.length ? (
        <div className="persons-table-wrap">
          <table className="persons-table">
            <thead>
              <tr>
                <th scope="col">ФИО</th>
                <th scope="col">Дата рождения</th>
                <th scope="col">Дата смерти</th>
                <th scope="col">Город</th>
              </tr>
            </thead>
            <tbody>
              {filteredItems.map((person) => (
                <tr key={person.id}>
                  <td>
                    <button
                      className="persons-table-name"
                      onClick={() =>
                        navigate(`/trees/${tree.id}?person=${person.id}`)
                      }
                    >
                      {fullName(person)}
                    </button>
                  </td>
                  <td>{formatDate(person.birth_date)}</td>
                  <td>{formatDate(person.death_date)}</td>
                  <td>{person.metadata?.city?.trim() || "Не указан"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : normalizedSearch ? (
        <EmptyState
          icon={<Search size={28} />}
          title="Ничего не найдено"
          text="Попробуйте изменить запрос поиска по ФИО."
        />
      ) : (
        <EmptyState
          icon={<ContactRound size={28} />}
          title="Персон пока нет"
          text="Добавьте первого человека в выбранное дерево."
          action={<Button onClick={() => navigate(`/trees/${tree.id}`)}>Открыть древо</Button>}
        />
      )}
    </div>
  );
}

function peopleCountLabel(count: number) {
  const remainder = count % 10;
  const lastTwoDigits = count % 100;
  if (remainder === 1 && lastTwoDigits !== 11) return "персона";
  if (
    remainder >= 2 &&
    remainder <= 4 &&
    (lastTwoDigits < 10 || lastTwoDigits >= 20)
  ) {
    return "персоны";
  }
  return "персон";
}
