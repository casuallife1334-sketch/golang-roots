import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Plus, Trash2, UserRound } from "lucide-react";
import { api } from "../../api";
import { useAuth } from "../../auth";
import { keys, relationshipsQuery } from "../../data/queries";
import { parentChild } from "../../graph/model";
import { Button, Field, Modal, Notice, Select } from "../../shared/ui";
import type {
  Family,
  FamilyMemberRole,
  Person,
  Relationship,
  Tree,
} from "../../types";
import { fullName, metadataComment } from "../../utils";
import { CommentEditor } from "../tree/CommentEditor";

const roles: { value: FamilyMemberRole; label: string }[] = [
  { value: "partner", label: "Партнёр" },
  { value: "parent", label: "Родитель" },
  { value: "child", label: "Ребёнок" },
];

type Confirm =
  | { kind: "family" }
  | { kind: "member"; id: string; label: string }
  | { kind: "relationship"; id: string; label: string };

export function FamilyDialog({
  tree,
  family,
  people,
  editable,
  onClose,
  onCreated,
  onDeleted,
  onOpenPerson,
}: {
  tree: Tree;
  family: Family | null;
  people: Person[];
  editable: boolean;
  onClose: () => void;
  onCreated: (id: string) => void;
  onDeleted: () => void;
  onOpenPerson: (id: string) => void;
}) {
  const { user } = useAuth();
  const query = useQueryClient();
  const [name, setName] = useState(family?.name ?? "");
  const [comment, setComment] = useState(metadataComment(family?.metadata));
  const [role, setRole] = useState<FamilyMemberRole>("parent");
  const [personId, setPersonId] = useState("");
  const [relationshipId, setRelationshipId] = useState("");
  const [confirm, setConfirm] = useState<Confirm | null>(null);
  const [error, setError] = useState("");
  const relationships = useQuery({
    ...relationshipsQuery(user?.id, tree.id),
    enabled: Boolean(family && user),
  });
  const peopleById = new Map(people.map((person) => [person.id, person]));
  const members = new Map<string, FamilyMemberRole[]>();
  for (const member of family?.members ?? []) {
    const current = members.get(member.person_id) ?? [];
    if (!current.includes(member.role)) current.push(member.role);
    members.set(member.person_id, current);
  }
  const linkedIds = new Set(family?.relationship_ids ?? []);
  const linked = (relationships.data ?? []).filter((item) =>
    linkedIds.has(item.id),
  );
  const availablePeople = people.filter(
    (person) => !members.get(person.id)?.includes(role),
  );
  const availableRelationships = (relationships.data ?? []).filter(
    (item) => !linkedIds.has(item.id),
  );
  const refresh = () =>
    query.invalidateQueries({ queryKey: keys.families(user?.id, tree.id) });
  const onError = (reason: Error) => setError(reason.message);

  const save = useMutation({
    mutationFn: () => {
      const metadata = { ...family?.metadata };
      if (comment.trim()) metadata.comment = comment.trim();
      else delete metadata.comment;
      const body = { name: name.trim(), metadata };
      return family
        ? api.patchFamily(tree.id, family.id, body)
        : api.createFamily(tree.id, body);
    },
    onSuccess: async (saved) => {
      await refresh();
      setError("");
      if (!family) onCreated(saved.id);
    },
    onError,
  });
  const changeMember = useMutation({
    mutationFn: (action: { kind: "add"; id: string; role: FamilyMemberRole } | { kind: "remove"; id: string }) =>
      action.kind === "add"
        ? api.addFamilyMember(tree.id, family!.id, {
            person_id: action.id,
            role: action.role,
          })
        : api.removeFamilyMember(tree.id, family!.id, action.id),
    onSuccess: async () => {
      setPersonId("");
      setConfirm(null);
      setError("");
      await refresh();
    },
    onError,
  });
  const changeRelationship = useMutation({
    mutationFn: (action: { kind: "attach" | "detach"; id: string }) =>
      action.kind === "attach"
        ? api.attachFamilyRelationship(tree.id, family!.id, action.id)
        : api.detachFamilyRelationship(tree.id, family!.id, action.id),
    onSuccess: async () => {
      setRelationshipId("");
      setConfirm(null);
      setError("");
      await refresh();
    },
    onError,
  });
  const removeFamily = useMutation({
    mutationFn: () => api.deleteFamily(tree.id, family!.id),
    onSuccess: async () => {
      await refresh();
      onDeleted();
    },
    onError,
  });
  const busy =
    save.isPending ||
    changeMember.isPending ||
    changeRelationship.isPending ||
    removeFamily.isPending;

  return (
    <>
      <Modal
        title={family ? "Семья" : "Создать семью"}
        onClose={onClose}
        busy={busy}
        wide={Boolean(family)}
        className="family-editor-modal"
      >
        <div className={`family-editor-layout ${family ? "with-members" : ""}`}>
          <form
            className="modal-form family-editor-form"
            onSubmit={(event) => {
              event.preventDefault();
              if (!busy && editable) {
                setError("");
                save.mutate();
              }
            }}
          >
            <fieldset disabled={!editable || busy} className="form-fields">
              <Field
                label="Название семьи"
                hint="Необязательно. Без названия показываются имена взрослых участников."
              >
                <input
                  autoFocus={editable}
                  maxLength={200}
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder="Например, семья Алексея и Анны"
                />
              </Field>
              {(editable || comment) && (
                <CommentEditor
                  title="Сведения о семье"
                  description="Необязательная заметка для этой семьи"
                  value={comment}
                  onChange={setComment}
                  placeholder="Добавьте сведения о семье…"
                  disabled={busy}
                  editable={editable}
                />
              )}
            </fieldset>
            {editable && (
              <div className={`modal-actions family-editor-actions ${family ? "with-delete" : ""}`}>
                {family && (
                  <Button
                    variant="danger"
                    disabled={busy}
                    onClick={() => {
                      setError("");
                      setConfirm({ kind: "family" });
                    }}
                  >
                    <Trash2 size={15} /> Удалить семью
                  </Button>
                )}
                <Button
                  type="submit"
                  loading={save.isPending}
                  disabled={busy && !save.isPending}
                >
                  {family ? "Сохранить изменения" : "Создать семью"}
                </Button>
              </div>
            )}
          </form>

          {family && (
            <div className="family-editor-management">
              <section className="family-editor-section">
                <div className="family-editor-section-heading">
                  <h3>Участники</h3>
                  <span>{members.size}</span>
                </div>
                {members.size ? (
                  <ul className="family-editor-items">
                    {[...members].map(([id, memberRoles]) => {
                      const person = peopleById.get(id);
                      const label = person ? fullName(person) : id;
                      return (
                        <li key={id} className="family-editor-item">
                          <div className="family-editor-item-text">
                            <strong>{label}</strong>
                            <small>
                              {memberRoles
                                .map((value) => roles.find((item) => item.value === value)?.label)
                                .join(" · ")}
                            </small>
                          </div>
                          {person && (
                            <button
                              type="button"
                              className="family-editor-icon"
                              aria-label={`Открыть ${label} в дереве`}
                              onClick={() => onOpenPerson(id)}
                            >
                              <UserRound size={16} />
                            </button>
                          )}
                          {editable && (
                            <button
                              type="button"
                              className="family-editor-icon danger"
                              disabled={busy}
                              aria-label={`Удалить ${label} из семьи`}
                              onClick={() => {
                                setError("");
                                setConfirm({ kind: "member", id, label });
                              }}
                            >
                              <Trash2 size={16} />
                            </button>
                          )}
                        </li>
                      );
                    })}
                  </ul>
                ) : (
                  <p className="family-editor-hint">В этой семье пока нет участников.</p>
                )}
                {editable && (
                  <div className="family-editor-add">
                    <div className="family-editor-pickers">
                      <Field label="Человек">
                        <Select
                          value={personId}
                          onChange={(event) => setPersonId(event.target.value)}
                          disabled={busy || !availablePeople.length}
                        >
                          <option value="">Выберите человека</option>
                          {availablePeople.map((person) => (
                            <option key={person.id} value={person.id}>
                              {fullName(person)}
                            </option>
                          ))}
                        </Select>
                      </Field>
                      <Field label="Роль">
                        <Select
                          value={role}
                          onChange={(event) => {
                            setRole(event.target.value as FamilyMemberRole);
                            setPersonId("");
                          }}
                          disabled={busy}
                        >
                          {roles.map((item) => (
                            <option key={item.value} value={item.value}>
                              {item.label}
                            </option>
                          ))}
                        </Select>
                      </Field>
                    </div>
                    <Button
                      variant="secondary"
                      disabled={busy || !personId || !availablePeople.some((person) => person.id === personId)}
                      loading={changeMember.isPending}
                      onClick={() => {
                        setError("");
                        changeMember.mutate({ kind: "add", id: personId, role });
                      }}
                    >
                      <Plus size={15} /> Добавить участника
                    </Button>
                  </div>
                )}
              </section>

              <section className="family-editor-section">
                <div className="family-editor-section-heading">
                  <h3>Привязанные связи</h3>
                  <span>{family.relationship_ids.length}</span>
                </div>
                {relationships.isPending ? (
                  <p className="family-editor-hint">Загрузка связей…</p>
                ) : relationships.error ? (
                  <Notice>{relationships.error.message}</Notice>
                ) : linked.length ? (
                  <ul className="family-editor-items">
                    {linked.map((relationship) => {
                      const label = relationshipLabel(relationship, peopleById);
                      return (
                        <li key={relationship.id} className="family-editor-item">
                          <div className="family-editor-item-text">
                            <strong>{label}</strong>
                          </div>
                          {editable && (
                            <button
                              type="button"
                              className="family-editor-icon danger"
                              disabled={busy}
                              aria-label={`Отвязать связь: ${label}`}
                              onClick={() => {
                                setError("");
                                setConfirm({ kind: "relationship", id: relationship.id, label });
                              }}
                            >
                              <Trash2 size={16} />
                            </button>
                          )}
                        </li>
                      );
                    })}
                  </ul>
                ) : (
                  <p className="family-editor-hint">Связи не привязаны к этой семье.</p>
                )}
                {editable && !relationships.error && !relationships.isPending && (
                  <div className="family-editor-add">
                    <Field label="Связь из дерева">
                      <Select
                        value={relationshipId}
                        onChange={(event) => setRelationshipId(event.target.value)}
                        disabled={busy || !availableRelationships.length}
                      >
                        <option value="">Выберите связь</option>
                        {availableRelationships.map((relationship) => (
                          <option key={relationship.id} value={relationship.id}>
                            {relationshipLabel(relationship, peopleById)}
                          </option>
                        ))}
                      </Select>
                    </Field>
                    <Button
                      variant="secondary"
                      disabled={busy || !relationshipId}
                      loading={changeRelationship.isPending}
                      onClick={() => {
                        setError("");
                        changeRelationship.mutate({ kind: "attach", id: relationshipId });
                      }}
                    >
                      <Plus size={15} /> Привязать связь
                    </Button>
                    <small className="family-editor-hint">
                      Привязка связи не добавляет участников автоматически.
                    </small>
                  </div>
                )}
              </section>
              <p className="family-editor-hint">
                Ручные изменения состава семьи сохраняются независимо от последующих
                изменений связей в дереве.
              </p>
            </div>
          )}
        </div>
        {error && <Notice>{error}</Notice>}
      </Modal>
      {confirm && (
        <Modal
          title={
            confirm.kind === "family"
              ? "Удалить семью?"
              : confirm.kind === "member"
                ? "Убрать участника?"
                : "Отвязать связь?"
          }
          onClose={() => setConfirm(null)}
          busy={busy}
        >
          <div className="confirm-content">
            <p>
              {confirm.kind === "family"
                ? "Запись семьи и её состав будут удалены. Люди и их связи в дереве останутся."
                : confirm.kind === "member"
                  ? `${confirm.label} исчезнет из участников семьи во всех ролях. Связи этого человека в дереве останутся.`
                  : `Связь «${confirm.label}» перестанет относиться к семье, но останется в дереве.`}
            </p>
            {error && <Notice>{error}</Notice>}
            <div className="modal-actions">
              <Button
                variant="secondary"
                disabled={busy}
                onClick={() => setConfirm(null)}
              >
                Отмена
              </Button>
              <Button
                variant="danger"
                loading={busy}
                onClick={() => {
                  setError("");
                  if (confirm.kind === "family") removeFamily.mutate();
                  else if (confirm.kind === "member")
                    changeMember.mutate({ kind: "remove", id: confirm.id });
                  else
                    changeRelationship.mutate({ kind: "detach", id: confirm.id });
                }}
              >
                {confirm.kind === "family"
                  ? "Удалить семью"
                  : confirm.kind === "member"
                    ? "Убрать участника"
                    : "Отвязать связь"}
              </Button>
            </div>
          </div>
        </Modal>
      )}
    </>
  );
}

function relationshipLabel(
  relationship: Relationship,
  people: Map<string, Person>,
) {
  const link = parentChild(relationship);
  if (link) {
    return `${fullName(people.get(link.parentId))} → ${fullName(people.get(link.childId))} (родитель и ребёнок)`;
  }
  return `${fullName(people.get(relationship.person1_id))} и ${fullName(people.get(relationship.person2_id))} (партнёры)`;
}
