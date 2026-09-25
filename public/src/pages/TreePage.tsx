import { useEffect, useState } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { Plus, Search, UserRound, X } from "lucide-react";
import { useAuth } from "../auth";
import { useTreeData } from "../data/queries";
import { FamilyGraph } from "../graph/FamilyGraph";
import { Button, EmptyState, Notice } from "../shared/ui";
import { ErrorBoundary } from "../shared/ErrorBoundary";
import { fullName } from "../utils";
import { TreeSwitcher } from "../features/tree/TreeSwitcher";
import { PersonPanel } from "../features/tree/PersonPanel";
import { PersonDialog } from "../features/tree/PersonDialog";
import { TreeDialog } from "../features/tree/TreeDialog";
import { RelationshipDialog } from "../features/tree/RelationshipDialog";

export function TreePage() {
  const { treeId } = useParams();
  // Local dialogs and selection cannot leak into another tree.
  return <TreeScreen key={treeId ?? "default"} treeId={treeId} />;
}
function TreeScreen({ treeId }: { treeId?: string }) {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [params, setParams] = useSearchParams();
  const { trees, tree, people, relationships } = useTreeData(treeId);
  const [selectedId, setSelectedId] = useState<string | null>(
    params.get("person"),
  );
  const [editingId, setEditingId] = useState<string | null>(null);
  const [dialog, setDialog] = useState<
    "tree" | "person" | "relationship" | null
  >(null);
  const [search, setSearch] = useState("");
  const selected =
    people.data?.find((person) => person.id === selectedId) ?? null;
  const editing =
    people.data?.find((person) => person.id === editingId) ?? null;
  const canWrite = tree?.role
    ? tree.role === "owner" || tree.role === "editor"
    : tree?.owner_id === user?.id;
  useEffect(() => {
    if (tree && user)
      sessionStorage.setItem(`roots:last-tree:${user.id}`, tree.id);
    if (!treeId && tree) navigate(`/trees/${tree.id}`, { replace: true });
  }, [tree?.id, treeId, user?.id, navigate]);
  useEffect(() => {
    setSelectedId(params.get("person"));
  }, [params]);
  const select = (id: string | null) => {
    setSelectedId(id);
    const next = new URLSearchParams(params);
    if (id) next.set("person", id);
    else next.delete("person");
    setParams(next, { replace: true });
  };
  const treeDialog = dialog === "tree" && (
    <TreeDialog
      onClose={() => setDialog(null)}
      onSaved={(id) => {
        setDialog(null);
        navigate(`/trees/${id}`);
      }}
    />
  );
  if (trees.isPending)
    return <div className="page-loader">Загрузка деревьев…</div>;
  if (trees.error)
    return (
      <div className="center-panel">
        <Notice>{trees.error.message}</Notice>
        <Button onClick={() => trees.refetch()}>Повторить</Button>
      </div>
    );
  if (!trees.data?.length)
    return (
      <>
        <EmptyState
          title="Создайте первое дерево"
          text="Начните с названия семейной истории."
          action={
            <Button onClick={() => setDialog("tree")}>Создать дерево</Button>
          }
        />
        {treeDialog}
      </>
    );
  if (!tree)
    return (
      <EmptyState
        title="Дерево не найдено"
        text="Оно удалено или недоступно вашему аккаунту."
        action={<Button onClick={() => navigate("/trees")}>К деревьям</Button>}
      />
    );
  const filtered = (people.data ?? []).filter((person) =>
    fullName(person).toLowerCase().includes(search.toLowerCase()),
  );
  return (
    <div className="tree-page">
      <header className="topbar">
        <div className="search-box">
          <Search size={18} />
          <input
            aria-label="Поиск по людям"
            placeholder="Поиск по имени…"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
          />
          {search && (
            <>
              <button aria-label="Очистить поиск" onClick={() => setSearch("")}>
                <X size={15} />
              </button>
              <div className="search-results">
                {filtered.slice(0, 10).map((person) => (
                  <button
                    key={person.id}
                    onClick={() => {
                      select(person.id);
                      setSearch("");
                    }}
                  >
                    {fullName(person)}
                  </button>
                ))}
                {!filtered.length && <span>Ничего не найдено</span>}
              </div>
            </>
          )}
        </div>
        <TreeSwitcher
          trees={trees.data}
          current={tree}
          onChange={(id) => navigate(`/trees/${id}`)}
          onCreate={() => setDialog("tree")}
        />
        {canWrite && (
          <Button
            onClick={() => {
              setEditingId(null);
              setDialog("person");
            }}
          >
            <Plus size={16} />
            Добавить человека
          </Button>
        )}
      </header>
      <div className="tree-content">
        <div className="canvas">
          {people.isPending || relationships.isPending ? (
            <div className="page-loader">Загрузка дерева…</div>
          ) : people.error || relationships.error ? (
            <div className="center-panel">
              <Notice>
                {people.error?.message || relationships.error?.message}
              </Notice>
              <Button
                onClick={() => {
                  void people.refetch();
                  void relationships.refetch();
                }}
              >
                Повторить
              </Button>
            </div>
          ) : people.data.length ? (
            <ErrorBoundary key={tree.id}>
              <div className="family-tree-library">
                <FamilyGraph
                  treeId={tree.id}
                  people={people.data}
                  relationships={relationships.data}
                  selectedId={selected?.id}
                  onSelect={(person) => select(person.id)}
                />
              </div>
            </ErrorBoundary>
          ) : (
            <EmptyState
              icon={<UserRound size={28} />}
              title="Дерево пока пусто"
              text="Добавьте первого человека."
              action={
                canWrite && (
                  <Button
                    onClick={() => {
                      setEditingId(null);
                      setDialog("person");
                    }}
                  >
                    Добавить первого человека
                  </Button>
                )
              }
            />
          )}
        </div>
        {selected && (
          <PersonPanel
            key={selected.id}
            person={selected}
            tree={tree}
            people={people.data ?? []}
            relationships={relationships.data ?? []}
            editable={canWrite}
            onClose={() => select(null)}
            onEdit={() => {
              setEditingId(selected.id);
              setDialog("person");
            }}
            onAddRelationship={() => setDialog("relationship")}
            onSelectPerson={(person) => select(person.id)}
          />
        )}
      </div>
      {dialog === "person" && canWrite && (
        <PersonDialog
          tree={tree}
          person={editing}
          onClose={() => setDialog(null)}
          onSaved={(person) => {
            select(person.id);
            setDialog(null);
          }}
        />
      )}
      {dialog === "relationship" && canWrite && (
        <RelationshipDialog
          tree={tree}
          people={people.data ?? []}
          relationships={relationships.data ?? []}
          selected={selected}
          onClose={() => setDialog(null)}
          onSaved={() => setDialog(null)}
        />
      )}
      {treeDialog}
    </div>
  );
}
