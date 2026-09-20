import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import "@xyflow/react/dist/style.css";
import {
  Plus,
  Search,
  ChevronDown,
  X,
  Camera,
  Upload,
  ZoomIn,
  Trash2,
  Link2,
  UserRound,
} from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { api, ApiError } from "../api";
import { useAuth } from "../auth";
import { usePersonPhotoUrls } from "../graph/photoUrls";
import { FamilyGraph } from "../graph/FamilyGraph";
import { Button, EmptyState, Field, Modal, Notice, Select } from "../shared/ui";
import type {
  Gender,
  Person,
  PersonInput,
  Relationship,
  RelationshipInput,
  Tree,
} from "../types";
import {
  dateInput,
  formatDate,
  formatYears,
  fullName,
  initials,
  validImage,
} from "../utils";

export function TreePage() {
  const { treeId } = useParams();
  const navigate = useNavigate();
  const query = useQueryClient();
  const { user } = useAuth();
  const [selected, setSelected] = useState<Person | null>(null);
  const [editingPerson, setEditingPerson] = useState<Person | null>(null);
  const [dialog, setDialog] = useState<
    "person" | "relationship" | "tree" | null
  >(null);
  const [search, setSearch] = useState("");
  const trees = useQuery({ queryKey: ["trees", user?.id], queryFn: api.trees });
  const selectedTree = treeId
    ? trees.data?.find((tree) => tree.id === treeId)
    : trees.data?.[0];
  const people = useQuery({
    queryKey: ["people", user?.id, selectedTree?.id],
    queryFn: () => api.persons(selectedTree!.id),
    enabled: Boolean(selectedTree?.id),
  });
  const relationships = useQuery({
    queryKey: ["relationships", user?.id, selectedTree?.id],
    queryFn: () => api.relationships(selectedTree!.id),
    enabled: Boolean(selectedTree?.id),
  });
  useEffect(() => {
    if (selectedTree && user)
      sessionStorage.setItem(`roots:last-tree:${user.id}`, selectedTree.id);
    if (!treeId && selectedTree)
      navigate(`/trees/${selectedTree.id}`, { replace: true });
    if (treeId && trees.data && !selectedTree)
      navigate("/trees", { replace: true });
    setSelected(null);
  }, [treeId, selectedTree?.id, trees.data]);
  const photoUrls = usePersonPhotoUrls(
    selectedTree?.id || "",
    people.data || [],
  );
  const filtered = (people.data || []).filter((person) =>
    `${person.first_name} ${person.last_name} ${String(person.metadata?.city || "")}`
      .toLowerCase()
      .includes(search.toLowerCase()),
  );
  if (trees.isLoading) return <PageLoading />;
  if (trees.isError)
    return (
      <ErrorPanel
        text="Не удалось загрузить деревья"
        retry={() => trees.refetch()}
      />
    );
  if (!trees.data?.length)
    return (
      <>
        <EmptyState
          icon={<LeafMark />}
          title="Создайте первое дерево"
          text="Начните с названия семейной истории, а затем добавьте первого человека."
          action={
            <Button onClick={() => setDialog("tree")}>
              <Plus size={16} />
              Создать дерево
            </Button>
          }
        />
        {dialog === "tree" && (
          <TreeDialog
            onClose={() => setDialog(null)}
            onSaved={(id) => {
              setDialog(null);
              query.invalidateQueries({ queryKey: ["trees"] });
              navigate(`/trees/${id}`);
            }}
          />
        )}
      </>
    );
  return (
    <div className="tree-page">
      <header className="topbar">
        <div className="search-box">
          <Search size={18} />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Поиск по имени или месту…"
            aria-label="Поиск по людям"
          />
          <kbd>⌘ K</kbd>
          {search && (
            <button onClick={() => setSearch("")} aria-label="Очистить">
              <X size={15} />
            </button>
          )}
          {search && (
            <div className="search-results">
              {filtered.length ? (
                filtered.slice(0, 6).map((person) => (
                  <button
                    key={person.id}
                    onClick={() => {
                      setSelected(person);
                      setSearch("");
                    }}
                  >
                    <span className="avatar tiny">{initials(person)}</span>
                    {fullName(person)}
                  </button>
                ))
              ) : (
                <span className="search-empty">Ничего не найдено</span>
              )}
            </div>
          )}
        </div>
        <TreeSwitcher
          trees={trees.data}
          current={selectedTree}
          onChange={(id) => navigate(`/trees/${id}`)}
          onCreate={() => setDialog("tree")}
        />
        <div className="top-actions">
          <Button
            onClick={() => {
              setEditingPerson(null);
              setDialog("person");
            }}
          >
            <Plus size={16} />
            Добавить человека
          </Button>
        </div>
      </header>
      <div className="tree-toolbar-mobile">
        <TreeSwitcher
          trees={trees.data}
          current={selectedTree}
          onChange={(id) => navigate(`/trees/${id}`)}
          onCreate={() => setDialog("tree")}
        />
      </div>
      <div className="tree-content">
        <div className="canvas">
          {people.data?.length ? (
            <div className="family-tree-library">
              <FamilyGraph
                people={people.data}
                relationships={relationships.data || []}
                photoUrls={photoUrls}
                onSelect={setSelected}
              />
            </div>
          ) : null}
          {!people.data?.length && (
            <EmptyState
              icon={<UserRound size={28} />}
              title="Дерево пока пусто"
              text="Добавьте первого человека, чтобы начать строить семейную историю."
              action={
                <Button onClick={() => setDialog("person")}>
                  <Plus size={16} />
                  Добавить первого человека
                </Button>
              }
            />
          )}
        </div>
        {selected && (
          <PersonPanel
            person={selected}
            people={people.data || []}
            relationships={relationships.data || []}
            tree={selectedTree!}
            owner={selectedTree!.owner_id === user?.id}
            onClose={() => setSelected(null)}
            onEdit={() => {
              setEditingPerson(selected);
              setDialog("person");
            }}
            onAddRelationship={() => setDialog("relationship")}
            onSelectPerson={setSelected}
            onChanged={() => {
              query.invalidateQueries({ queryKey: ["people"] });
              query.invalidateQueries({ queryKey: ["relationships"] });
            }}
          />
        )}
      </div>
      {dialog === "person" && (
        <PersonDialog
          tree={selectedTree!}
          person={editingPerson}
          onClose={() => {
            setEditingPerson(null);
            setDialog(null);
          }}
          onSaved={(person) => {
            setSelected(person);
            setEditingPerson(null);
            setDialog(null);
            query.invalidateQueries({ queryKey: ["people"] });
          }}
        />
      )}
      {dialog === "relationship" && (
        <RelationshipDialog
          tree={selectedTree!}
          people={people.data || []}
          selected={selected}
          relationships={relationships.data || []}
          onClose={() => setDialog(null)}
          onSaved={() => {
            setDialog(null);
            query.invalidateQueries({ queryKey: ["relationships"] });
          }}
        />
      )}
      {dialog === "tree" && (
        <TreeDialog
          onClose={() => setDialog(null)}
          onSaved={(id) => {
            setDialog(null);
            query.invalidateQueries({ queryKey: ["trees"] });
            navigate(`/trees/${id}`);
          }}
        />
      )}
    </div>
  );
}
function LeafMark() {
  return (
    <span className="empty-leaf">
      <Plus size={23} />
    </span>
  );
}
function PageLoading() {
  return (
    <div className="page-loader">
      <span className="spinner" />
      Загрузка дерева
    </div>
  );
}
function ErrorPanel({ text, retry }: { text: string; retry: () => void }) {
  return (
    <div className="center-panel">
      <Notice>{text}</Notice>
      <Button variant="secondary" onClick={retry}>
        Повторить
      </Button>
    </div>
  );
}
function TreeSwitcher({
  trees,
  current,
  onChange,
  onCreate,
}: {
  trees: Tree[];
  current?: Tree;
  onChange: (id: string) => void;
  onCreate: () => void;
}) {
  const [open, setOpen] = useState(false);
  return (
    <div className={`tree-switcher ${open ? "open" : ""}`}>
      <button
        className="tree-switcher-trigger"
        onClick={() => setOpen(!open)}
        aria-expanded={open}
        aria-label="Выбрать дерево"
      >
        <span>{current?.name || "Выберите дерево"}</span>
        <ChevronDown size={15} />
      </button>
      {open && (
        <div className="tree-menu" role="menu">
          {trees.map((tree) => (
            <button
              key={tree.id}
              className={tree.id === current?.id ? "selected" : ""}
              onClick={() => {
                onChange(tree.id);
                setOpen(false);
              }}
              role="menuitem"
            >
              {tree.name}
              {tree.id === current?.id && <span>Текущее</span>}
            </button>
          ))}
        </div>
      )}
      <button
        onClick={onCreate}
        title="Создать дерево"
        aria-label="Создать дерево"
      >
        <Plus size={15} />
      </button>
    </div>
  );
}
function PersonPanel({
  person,
  people,
  relationships,
  tree,
  owner,
  onClose,
  onEdit,
  onAddRelationship,
  onSelectPerson,
  onChanged,
}: {
  person: Person;
  people: Person[];
  relationships: Relationship[];
  tree: Tree;
  owner: boolean;
  onClose: () => void;
  onEdit: () => void;
  onAddRelationship: () => void;
  onSelectPerson: (person: Person) => void;
  onChanged: () => void;
}) {
  const query = useQueryClient();
  const [tab, setTab] = useState<"overview" | "relations">("overview");
  const [photo, setPhoto] = useState<string>();
  const [cropFile, setCropFile] = useState<File>();
  const [photoError, setPhotoError] = useState("");
  const [confirmDelete, setConfirmDelete] = useState<"person" | string | null>(
    null,
  );
  const file = useRef<HTMLInputElement>(null);
  const deletePerson = useMutation({
    mutationFn: () => api.deletePerson(tree.id, person.id),
    onSuccess: () => {
      query.invalidateQueries({ queryKey: ["people"] });
      query.invalidateQueries({ queryKey: ["relationships"] });
      onClose();
    },
  });
  const personRelationships = relationships.filter(
    (r) => r.person1_id === person.id || r.person2_id === person.id,
  );
  const connected = personRelationships
    .map((r) =>
      people.find(
        (p) =>
          p.id === (r.person1_id === person.id ? r.person2_id : r.person1_id),
      ),
    )
    .filter(Boolean) as Person[];
  useEffect(() => {
    if (!person.photo_url) {
      setPhoto(undefined);
      return;
    }
    let active = true;
    let objectUrl: string | undefined;
    api
      .photo(tree.id, person.id)
      .then((blob) => {
        if (!active) return;
        objectUrl = URL.createObjectURL(blob);
        setPhoto(objectUrl);
      })
      .catch(() => setPhoto(undefined));
    return () => {
      active = false;
      if (objectUrl) URL.revokeObjectURL(objectUrl);
    };
  }, [person.id, person.photo_url, tree.id]);
  const upload = async (selected: File) => {
    if (!validImage(selected)) {
      setPhotoError("Нужен JPEG, PNG или WebP до 9,5 МБ");
      return;
    }
    try {
      setPhotoError("");
      await api.uploadPhoto(tree.id, person.id, selected);
      query.invalidateQueries({ queryKey: ["people"] });
      query.invalidateQueries({ queryKey: ["photo", tree.id, person.id] });
      onChanged();
    } catch (error) {
      setPhotoError(
        error instanceof ApiError
          ? error.message
          : "Не удалось загрузить фотографию",
      );
    }
  };
  const removePhoto = async () => {
    try {
      setPhotoError("");
      await api.deletePhoto(tree.id, person.id);
      setPhoto(undefined);
      query.invalidateQueries({ queryKey: ["people"] });
      query.invalidateQueries({ queryKey: ["photo", tree.id, person.id] });
      onChanged();
    } catch (error) {
      setPhotoError(
        error instanceof ApiError
          ? error.message
          : "Не удалось удалить фотографию",
      );
    }
  };
  return (
    <>
      <aside className="detail-panel">
        <div className="detail-head">
          <span>Человек</span>
          <button
            className="icon-button"
            onClick={onClose}
            aria-label="Закрыть"
          >
            <X size={18} />
          </button>
        </div>
        <div className="profile-summary">
          <div className="profile-photo">
            {photo ? (
              <img src={photo} alt="" />
            ) : (
              <span>{initials(person)}</span>
            )}
            <button
              className="camera-button"
              onClick={() => file.current?.click()}
              aria-label="Загрузить фотографию"
            >
              <Camera size={15} />
            </button>
            {person.photo_url && owner && (
              <button
                className="camera-button remove-photo"
                onClick={removePhoto}
                aria-label="Удалить фотографию"
                title="Удалить фотографию"
              >
                <Trash2 size={14} />
              </button>
            )}
            <input
              ref={file}
              hidden
              type="file"
              accept="image/jpeg,image/png,image/webp"
              onChange={(e) => {
                const next = e.target.files?.[0];
                if (next) setCropFile(next);
              }}
            />
          </div>
          <div className="profile-meta">
            <h2>{fullName(person)}</h2>
            <p>
              {person.birth_date
                ? formatDate(person.birth_date)
                : "Дата рождения не указана"}
            </p>
            {person.metadata?.city && (
              <p className="profile-location">{String(person.metadata.city)}</p>
            )}
            {!person.death_date && <span className="life-status">Жив(а)</span>}
            {owner && (
              <div className="profile-actions">
                <Button variant="secondary" onClick={onEdit}>
                  Редактировать
                </Button>
                <Button
                  variant="danger"
                  onClick={() => setConfirmDelete("person")}
                >
                  Удалить
                </Button>
              </div>
            )}
          </div>
        </div>
        {photoError && (
          <div className="profile-error">
            <Notice>{photoError}</Notice>
          </div>
        )}
        <div className="tabs">
          <button
            className={tab === "overview" ? "active" : ""}
            onClick={() => setTab("overview")}
          >
            Обзор
          </button>
          <button
            className={tab === "relations" ? "active" : ""}
            onClick={() => setTab("relations")}
          >
            Родственники
          </button>
        </div>
        {tab === "overview" ? (
          <div className="detail-scroll">
            <h3>Основная информация</h3>
            <InfoRow
              label="Дата смерти"
              value={
                person.death_date ? formatDate(person.death_date) : "Не указана"
              }
            />
            <InfoRow
              label="Пол"
              value={
                person.gender === "male"
                  ? "Мужской"
                  : person.gender === "female"
                    ? "Женский"
                    : "Не указан"
              }
            />
            {person.metadata?.occupation && (
              <InfoRow
                label="Род деятельности"
                value={String(person.metadata.occupation)}
              />
            )}
          </div>
        ) : (
          <div className="detail-scroll">
            <h3>Связи человека</h3>
            {personRelationships.length ? (
              personRelationships.map((relationship) => {
                const item = people.find(
                  (candidate) =>
                    candidate.id ===
                    (relationship.person1_id === person.id
                      ? relationship.person2_id
                      : relationship.person1_id),
                );
                if (!item) return null;
                return (
                  <div className="relation-row" key={relationship.id}>
                    <button
                      className="relation-main"
                      onClick={() => onSelectPerson(item)}
                    >
                      <span className="avatar tiny">{initials(item)}</span>
                      <span>{fullName(item)}</span>
                      <small>
                        {relationship.type === "spouse"
                          ? "Партнёр"
                          : "Родственник"}
                      </small>
                    </button>
                    {owner && (
                      <button
                        className="relation-delete"
                        onClick={() => setConfirmDelete(relationship.id)}
                        aria-label="Удалить связь"
                        title="Удалить связь"
                      >
                        <Trash2 size={14} />
                      </button>
                    )}
                  </div>
                );
              })
            ) : (
              <p className="muted">Связей пока нет.</p>
            )}
            {owner && (
              <Button variant="secondary" onClick={onAddRelationship}>
                <Link2 size={15} />
                Добавить связь
              </Button>
            )}
          </div>
        )}
      </aside>
      {cropFile && (
        <PhotoCropModal
          file={cropFile}
          onCancel={() => setCropFile(undefined)}
          onConfirm={(cropped) => {
            setCropFile(undefined);
            void upload(cropped);
          }}
        />
      )}
      {confirmDelete && (
        <Modal
          title={
            confirmDelete === "person" ? "Удалить человека?" : "Удалить связь?"
          }
          onClose={() => setConfirmDelete(null)}
        >
          <div className="confirm-content">
            <Notice>
              {confirmDelete === "person"
                ? "Будут удалены человек и его связи."
                : "Люди останутся в дереве, удалится только эта связь."}
            </Notice>
            <div className="modal-actions">
              <Button
                variant="secondary"
                onClick={() => setConfirmDelete(null)}
              >
                Отмена
              </Button>
              <Button
                onClick={() => {
                  if (confirmDelete === "person") deletePerson.mutate();
                  else
                    api
                      .deleteRelationship(tree.id, confirmDelete)
                      .then(onChanged);
                  setConfirmDelete(null);
                }}
              >
                Удалить
              </Button>
            </div>
          </div>
        </Modal>
      )}
    </>
  );
}
function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="info-row">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}
function PersonDialog({
  tree,
  person,
  onClose,
  onSaved,
}: {
  tree: Tree;
  person: Person | null;
  onClose: () => void;
  onSaved: (person: Person) => void;
}) {
  const query = useQueryClient();
  const isEdit = Boolean(person);
  const [form, setForm] = useState({
    first_name: person?.first_name || "",
    patronymic: String(person?.metadata?.patronymic || ""),
    last_name: person?.last_name || "",
    birth_date: dateInput(person?.birth_date),
    death_date: dateInput(person?.death_date),
    gender: person?.gender || "",
    city: String(person?.metadata?.city || ""),
    occupation: String(person?.metadata?.occupation || ""),
  });
  const [file, setFile] = useState<File>();
  const [cropFile, setCropFile] = useState<File>();
  const [removeExistingPhoto, setRemoveExistingPhoto] = useState(false);
  const [error, setError] = useState("");
  const mutation = useMutation({
    mutationFn: async () => {
      const metadata = {
        ...(person?.metadata || {}),
        ...(form.patronymic ? { patronymic: form.patronymic.trim() } : {}),
        ...(form.city ? { city: form.city } : {}),
        ...(form.occupation ? { occupation: form.occupation } : {}),
      };
      const body: PersonInput = {
        first_name: form.first_name.trim(),
        last_name: form.last_name.trim(),
        birth_date: form.birth_date ? `${form.birth_date}T00:00:00Z` : null,
        death_date: form.death_date ? `${form.death_date}T00:00:00Z` : null,
        gender: (form.gender || null) as Gender | null,
        metadata,
      };
      const saved = isEdit
        ? await api.patchPerson(tree.id, person!.id, {
            first_name: body.first_name,
            last_name: body.last_name,
            birth_date: body.birth_date,
            death_date: body.death_date,
            gender: body.gender,
            metadata,
          })
        : await api.createPerson(tree.id, body);
      if (file) await api.uploadPhoto(tree.id, saved.id, file);
      else if (isEdit && removeExistingPhoto && person?.photo_url)
        await api.deletePhoto(tree.id, saved.id);
      return isEdit ? { ...saved } : saved;
    },
    onSuccess: (saved) => {
      query.invalidateQueries({ queryKey: ["people"] });
      onSaved(saved);
    },
    onError: (e) =>
      setError(
        e instanceof ApiError ? e.message : "Не удалось сохранить человека",
      ),
  });
  const update = (key: keyof typeof form, value: string) =>
    setForm((current) => ({ ...current, [key]: value }));
  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (!form.first_name.trim() || !form.last_name.trim())
      return setError("Укажите имя и фамилию");
    if (form.birth_date && form.death_date && form.birth_date > form.death_date)
      return setError("Дата рождения не может быть позже даты смерти");
    if (file && !validImage(file))
      return setError("Поддерживаются JPEG, PNG и WebP до 9,5 МБ");
    mutation.mutate();
  };
  return (
    <>
      <Modal
        title={isEdit ? "Редактировать человека" : "Добавить человека"}
        onClose={onClose}
      >
        <form className="modal-form" onSubmit={submit}>
          {error && <Notice>{error}</Notice>}
          <div className="form-grid">
            <Field label="Имя">
              <input
                value={form.first_name}
                onChange={(e) => update("first_name", e.target.value)}
                autoFocus
                required
              />
            </Field>
            <Field label="Отчество">
              <input
                value={form.patronymic}
                onChange={(e) => update("patronymic", e.target.value)}
              />
            </Field>
            <Field label="Фамилия">
              <input
                value={form.last_name}
                onChange={(e) => update("last_name", e.target.value)}
                required
              />
            </Field>
            <Field label="Дата рождения">
              <input
                type="date"
                value={form.birth_date}
                onChange={(e) => update("birth_date", e.target.value)}
              />
            </Field>
            <Field label="Дата смерти">
              <input
                type="date"
                value={form.death_date}
                onChange={(e) => update("death_date", e.target.value)}
              />
            </Field>
          </div>
          <Field label="Пол">
            <Select
              value={form.gender}
              onChange={(e) => update("gender", e.target.value)}
            >
              <option value="">Не указан</option>
              <option value="female">Женский</option>
              <option value="male">Мужской</option>
              <option value="other">Другой</option>
            </Select>
          </Field>
          <div className="form-grid">
            <Field label="Место / город">
              <input
                value={form.city}
                onChange={(e) => update("city", e.target.value)}
              />
            </Field>
            <Field label="Род деятельности">
              <input
                value={form.occupation}
                onChange={(e) => update("occupation", e.target.value)}
              />
            </Field>
          </div>
          <Field label="Фотография">
            <label className="file-picker">
              <Upload size={17} />
              <span>
                {file?.name ||
                  (removeExistingPhoto
                    ? "Фотография будет удалена"
                    : person?.photo_url
                      ? "Заменить фотографию"
                      : "Добавить фотографию")}
              </span>
              <input
                type="file"
                accept="image/jpeg,image/png,image/webp"
                onChange={(e) => {
                  const next = e.target.files?.[0];
                  if (next) setCropFile(next);
                  setRemoveExistingPhoto(false);
                }}
              />
              {(file || person?.photo_url) && (
                <button
                  type="button"
                  className="file-remove"
                  onClick={(event) => {
                    event.preventDefault();
                    setFile(undefined);
                    if (person?.photo_url) setRemoveExistingPhoto(true);
                  }}
                  aria-label="Удалить фотографию"
                  title="Удалить фотографию"
                >
                  <Trash2 size={14} />
                </button>
              )}
            </label>
          </Field>
          <div className="modal-actions">
            <Button variant="secondary" type="button" onClick={onClose}>
              Отмена
            </Button>
            <Button type="submit" loading={mutation.isPending}>
              {isEdit ? "Сохранить" : "Добавить человека"}
            </Button>
          </div>
        </form>
      </Modal>
      {cropFile && (
        <PhotoCropModal
          file={cropFile}
          onCancel={() => setCropFile(undefined)}
          onConfirm={(cropped) => {
            setFile(cropped);
            setCropFile(undefined);
          }}
        />
      )}
    </>
  );
}
function PhotoCropModal({
  file,
  onCancel,
  onConfirm,
}: {
  file: File;
  onCancel: () => void;
  onConfirm: (file: File) => void;
}) {
  const [source, setSource] = useState<string>();
  const [zoom, setZoom] = useState(1);
  const [panX, setPanX] = useState(0);
  const [panY, setPanY] = useState(0);
  const [loading, setLoading] = useState(false);
  const drag = useRef<{
    startX: number;
    startY: number;
    panX: number;
    panY: number;
  } | null>(null);
  useEffect(() => {
    const url = URL.createObjectURL(file);
    setSource(url);
    return () => URL.revokeObjectURL(url);
  }, [file]);
  const crop = () => {
    if (!source) return;
    setLoading(true);
    const image = new Image();
    image.onload = () => {
      const width = 420;
      const height = 600;
      const scale =
        Math.max(width / image.naturalWidth, height / image.naturalHeight) *
        zoom;
      const canvas = document.createElement("canvas");
      canvas.width = width;
      canvas.height = height;
      const context = canvas.getContext("2d");
      if (!context) return;
      context.drawImage(
        image,
        (width - image.naturalWidth * scale) / 2 + panX,
        (height - image.naturalHeight * scale) / 2 + panY,
        image.naturalWidth * scale,
        image.naturalHeight * scale,
      );
      canvas.toBlob(
        (blob) => {
          if (blob)
            onConfirm(
              new File([blob], file.name.replace(/\.[^.]+$/, ".jpg"), {
                type: "image/jpeg",
              }),
            );
          setLoading(false);
        },
        "image/jpeg",
        0.92,
      );
    };
    image.src = source;
  };
  return (
    <Modal title="Настроить фотографию" onClose={onCancel}>
      <div className="crop-form">
        <div className="crop-preview">
          <div
            className="crop-frame"
            onPointerDown={(event) => {
              event.currentTarget.setPointerCapture(event.pointerId);
              drag.current = {
                startX: event.clientX,
                startY: event.clientY,
                panX,
                panY,
              };
            }}
            onPointerMove={(event) => {
              if (!drag.current) return;
              setPanX(
                Math.max(
                  -120,
                  Math.min(
                    120,
                    drag.current.panX + (event.clientX - drag.current.startX),
                  ),
                ),
              );
              setPanY(
                Math.max(
                  -150,
                  Math.min(
                    150,
                    drag.current.panY + (event.clientY - drag.current.startY),
                  ),
                ),
              );
            }}
            onPointerUp={() => {
              drag.current = null;
            }}
            onPointerCancel={() => {
              drag.current = null;
            }}
          >
            {source && (
              <img
                src={source}
                alt="Предпросмотр фотографии"
                style={{
                  transform: `translate(calc(-50% + ${(panX * 2) / 3}px), calc(-50% + ${(panY * 2) / 3}px)) scale(${zoom})`,
                }}
              />
            )}
          </div>
          <span className="crop-ratio">21 × 30</span>
        </div>
        <div className="crop-controls">
          <p className="crop-help">
            Перетащите фотографию внутри рамки, чтобы выбрать кадрирование
          </p>
          <label>
            <span>
              <ZoomIn size={15} />
              Масштаб
            </span>
            <input
              type="range"
              min="1"
              max="2.5"
              step=".01"
              value={zoom}
              onChange={(e) => setZoom(Number(e.target.value))}
            />
          </label>
        </div>
        <div className="modal-actions">
          <Button variant="secondary" onClick={onCancel}>
            Отмена
          </Button>
          <Button loading={loading} onClick={crop}>
            Применить фото
          </Button>
        </div>
      </div>
    </Modal>
  );
}
function RelationshipDialog({
  tree,
  people,
  selected,
  relationships,
  onClose,
  onSaved,
}: {
  tree: Tree;
  people: Person[];
  selected: Person | null;
  relationships: Relationship[];
  onClose: () => void;
  onSaved: () => void;
}) {
  const [type, setType] = useState<"parent_child" | "spouse">("parent_child");
  const [one, setOne] = useState(selected?.id || "");
  const [two, setTwo] = useState("");
  const [direction, setDirection] = useState<"parent" | "child">("parent");
  const [error, setError] = useState("");
  const mutation = useMutation({
    mutationFn: (body: RelationshipInput) =>
      api.createRelationship(tree.id, body),
    onSuccess: onSaved,
    onError: (e) =>
      setError(e instanceof ApiError ? e.message : "Не удалось создать связь"),
  });
  const choices = people.filter((person) => person.id !== one);
  const duplicate = relationships.some((r) => {
    const samePair =
      new Set([r.person1_id, r.person2_id]).size === new Set([one, two]).size &&
      r.person1_id !== r.person2_id &&
      [r.person1_id, r.person2_id].every((id) => [one, two].includes(id));
    return samePair && r.type === type;
  });
  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!one || !two || one === two)
      return setError("Выберите двух разных людей");
    if (duplicate) return setError("Такая связь уже существует");
    mutation.mutate({
      person1_id: one,
      person2_id: two,
      type,
      ...(type === "parent_child" ? { direction } : {}),
    });
  };
  return (
    <Modal title="Добавить связь" onClose={onClose}>
      <form className="modal-form" onSubmit={submit}>
        {error && <Notice>{error}</Notice>}
        <Field label="Тип связи">
          <Select
            value={type}
            onChange={(e) =>
              setType(e.target.value as "parent_child" | "spouse")
            }
          >
            <option value="parent_child">Родитель → ребёнок</option>
            <option value="spouse">Партнёры / супруги</option>
          </Select>
        </Field>
        <Field
          label={
            type === "spouse"
              ? "Первый человек"
              : direction === "parent"
                ? "Родитель"
                : "Ребёнок"
          }
        >
          <Select value={one} onChange={(e) => setOne(e.target.value)}>
            <option value="">Выберите человека</option>
            {people.map((p) => (
              <option key={p.id} value={p.id}>
                {fullName(p)}
              </option>
            ))}
          </Select>
        </Field>
        {type === "parent_child" && (
          <Field label="Направление">
            <Select
              value={direction}
              onChange={(e) =>
                setDirection(e.target.value as "parent" | "child")
              }
            >
              <option value="parent">Родитель → ребёнок</option>
              <option value="child">Ребёнок → родитель</option>
            </Select>
          </Field>
        )}
        <Field
          label={
            type === "spouse"
              ? "Второй человек"
              : direction === "parent"
                ? "Ребёнок"
                : "Родитель"
          }
        >
          <Select value={two} onChange={(e) => setTwo(e.target.value)}>
            <option value="">Выберите человека</option>
            {choices.map((p) => (
              <option key={p.id} value={p.id}>
                {fullName(p)}
              </option>
            ))}
          </Select>
        </Field>
        <div className="modal-actions">
          <Button variant="secondary" type="button" onClick={onClose}>
            Отмена
          </Button>
          <Button type="submit" loading={mutation.isPending}>
            Добавить связь
          </Button>
        </div>
      </form>
    </Modal>
  );
}
function TreeDialog({
  onClose,
  onSaved,
}: {
  onClose: () => void;
  onSaved: (id: string) => void;
}) {
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const mutation = useMutation({
    mutationFn: () => api.createTree(name.trim()),
    onSuccess: (tree) => onSaved(tree.id),
    onError: (e) =>
      setError(e instanceof ApiError ? e.message : "Не удалось создать дерево"),
  });
  return (
    <Modal title="Новое дерево" onClose={onClose}>
      <form
        className="modal-form"
        onSubmit={(e) => {
          e.preventDefault();
          if (!name.trim()) return setError("Введите название дерева");
          mutation.mutate();
        }}
      >
        {error && <Notice>{error}</Notice>}
        <Field label="Название дерева">
          <input
            autoFocus
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Семья Ивановых"
          />
        </Field>
        <div className="modal-actions">
          <Button variant="secondary" type="button" onClick={onClose}>
            Отмена
          </Button>
          <Button type="submit" loading={mutation.isPending}>
            Создать дерево
          </Button>
        </div>
      </form>
    </Modal>
  );
}
