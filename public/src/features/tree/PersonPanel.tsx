import { useRef, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import {
  Camera,
  Link2,
  SquareArrowOutUpRight,
  Trash2,
  X,
} from "lucide-react";
import { api } from "../../api";
import type { Person, Relationship, Tree } from "../../types";
import { useTreeCache } from "../../data/queries";
import { usePersonPhoto } from "../../data/photos";
import {
  fullName,
  initials,
  formatDate,
  metadataComment,
  validImage,
} from "../../utils";
import { Button, Modal, Notice } from "../../shared/ui";
import { PhotoCropDialog } from "./PhotoCropDialog";
import { parentChild } from "../../graph/model";
import { RelationshipDetailsDialog } from "./RelationshipDetailsDialog";
import { DocumentsSection } from "./DocumentsSection";

export function PersonPanel({
  person,
  people,
  relationships,
  tree,
  editable,
  onClose,
  onEdit,
  onAddRelationship,
  onSelectPerson,
}: {
  person: Person;
  people: Person[];
  relationships: Relationship[];
  tree: Tree;
  editable: boolean;
  onClose: () => void;
  onEdit: () => void;
  onAddRelationship: () => void;
  onSelectPerson: (person: Person) => void;
}) {
  const cache = useTreeCache(tree.id);
  const { url: photo } = usePersonPhoto(tree.id, person);
  const file = useRef<HTMLInputElement>(null);
  const [tab, setTab] = useState("overview");
  const [cropFile, setCropFile] = useState<File>();
  const [error, setError] = useState("");
  const [confirm, setConfirm] = useState<string | null>(null);
  const [commentRelationship, setCommentRelationship] =
    useState<Relationship | null>(null);
  const remove = useMutation({
    mutationFn: (id: string) =>
      id === "person"
        ? api.deletePerson(tree.id, person.id)
        : api.deleteRelationship(tree.id, id),
    onSuccess: async (_, id) => {
      await cache.refresh();
      setConfirm(null);
      if (id === "person") onClose();
    },
  });
  const photoMutation = useMutation({
    mutationFn: async (next: File | null) => {
      if (next)
        cache.savePerson(await api.uploadPhoto(tree.id, person.id, next));
      else {
        await api.deletePhoto(tree.id, person.id);
        cache.savePerson({ ...person, photo_url: null });
      }
    },
    onSuccess: async () => {
      await cache.refreshPhoto(person.id);
      await cache.refresh();
    },
  });
  const related = relationships.filter(
    (item) => item.person1_id === person.id || item.person2_id === person.id,
  );
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
            {editable && (
              <>
                <button
                  className="camera-button"
                  disabled={photoMutation.isPending}
                  onClick={() => file.current?.click()}
                  aria-label="Загрузить фотографию"
                >
                  <Camera size={15} />
                </button>
                {person.photo_url && (
                  <button
                    className="camera-button remove-photo"
                    disabled={photoMutation.isPending}
                    onClick={() => photoMutation.mutate(null)}
                    aria-label="Удалить фотографию"
                  >
                    <Trash2 size={14} />
                  </button>
                )}
                <input
                  ref={file}
                  hidden
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  onChange={(event) => {
                    const next = event.target.files?.[0];
                    event.target.value = "";
                    setError("");
                    if (!next) return;
                    if (!validImage(next))
                      return setError("Нужен JPEG, PNG или WebP до 9,5 МБ");
                    setCropFile(next);
                  }}
                />
              </>
            )}
          </div>
          <div className="profile-meta">
            <h2>{fullName(person)}</h2>
            <p>
              {person.birth_date
                ? formatDate(person.birth_date)
                : "Дата рождения не указана"}
            </p>
            {person.metadata?.city && <p>{person.metadata.city}</p>}
            {!person.death_date && <span className="life-status">Жив(а)</span>}
            {editable && (
              <div className="profile-actions">
                <Button variant="secondary" onClick={onEdit}>
                  Редактировать
                </Button>
                <Button
                  variant="danger"
                  onClick={() => {
                    remove.reset();
                    setConfirm("person");
                  }}
                >
                  Удалить
                </Button>
              </div>
            )}
          </div>
        </div>
        {(error || photoMutation.error) && (
          <Notice>{error || photoMutation.error?.message}</Notice>
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
          <button
            className={tab === "documents" ? "active" : ""}
            onClick={() => setTab("documents")}
          >
            Документы
          </button>
        </div>
        <div className="detail-scroll">
          {tab === "overview" ? (
            <>
              <h3>Основная информация</h3>
              <InfoRow
                label="Дата смерти"
                value={
                  person.death_date
                    ? formatDate(person.death_date)
                    : "Не указана"
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
                  value={person.metadata.occupation}
                />
              )}
              {person.metadata?.comment && (
                <>
                  <h3 className="section-gap">Комментарий</h3>
                  <p className="person-comment-view">
                    {person.metadata.comment}
                  </p>
                </>
              )}
            </>
          ) : tab === "relations" ? (
            <>
              <h3>Связи человека</h3>
              {related.map((item) => {
                const relative = people.find(
                  (value) =>
                    value.id ===
                    (item.person1_id === person.id
                      ? item.person2_id
                      : item.person1_id),
                );
                if (!relative) return null;
                const comment = metadataComment(item.metadata);
                return (
                  <div
                    className={`relation-row ${comment ? "has-comment" : "no-comment"}`}
                    key={item.id}
                  >
                    <div className="relation-row-main">
                      <button
                        className="relation-main"
                        onClick={() => onSelectPerson(relative)}
                      >
                        <span className="avatar tiny">{initials(relative)}</span>
                        <span>{fullName(relative)}</span>
                        <small>
                          {item.type === "spouse"
                            ? "Партнёр"
                            : parentChild(item)?.parentId === person.id
                              ? "Ребёнок"
                              : "Родитель"}
                        </small>
                      </button>
                      {(editable || comment) && (
                        <button
                          className="relation-comment-button"
                          aria-label={
                            comment
                              ? "Открыть сведения о связи"
                              : "Добавить сведения о связи"
                          }
                          onClick={() => setCommentRelationship(item)}
                        >
                          <SquareArrowOutUpRight size={14} />
                        </button>
                      )}
                      {editable && (
                        <button
                          className="relation-delete"
                          aria-label="Удалить связь"
                          onClick={() => {
                            remove.reset();
                            setConfirm(item.id);
                          }}
                        >
                          <Trash2 size={14} />
                        </button>
                      )}
                    </div>
                    {comment && (
                      <button
                        className="relation-comment-preview"
                        onClick={() => setCommentRelationship(item)}
                      >
                        {comment}
                      </button>
                    )}
                  </div>
                );
              })}
              {!related.length && <p>Связей пока нет.</p>}
              {editable && (
                <Button variant="secondary" onClick={onAddRelationship}>
                  <Link2 size={15} />
                  Добавить связь
                </Button>
              )}
            </>
          ) : (
            <DocumentsSection
              treeId={tree.id}
              ownerType="person"
              ownerId={person.id}
              editable={editable}
            />
          )}
        </div>
      </aside>
      {cropFile && (
        <PhotoCropDialog
          file={cropFile}
          onCancel={() => setCropFile(undefined)}
          onConfirm={(value) => {
            setCropFile(undefined);
            photoMutation.mutate(value);
          }}
        />
      )}
      {confirm && (
        <Modal
          title={confirm === "person" ? "Удалить человека?" : "Удалить связь?"}
          onClose={() => setConfirm(null)}
          busy={remove.isPending}
        >
          <div className="confirm-content">
            <p>
              {confirm === "person"
                ? "Будут удалены человек и его связи."
                : "Люди останутся в дереве, удалится только эта связь."}
            </p>
            {remove.error && <Notice>{remove.error.message}</Notice>}
            <div className="modal-actions">
              <Button
                variant="secondary"
                disabled={remove.isPending}
                onClick={() => setConfirm(null)}
              >
                Отмена
              </Button>
              <Button
                variant="danger"
                loading={remove.isPending}
                onClick={() => remove.mutate(confirm)}
              >
                Удалить
              </Button>
            </div>
          </div>
        </Modal>
      )}
      {commentRelationship && (
        <RelationshipDetailsDialog
          tree={tree}
          relationship={commentRelationship}
          people={people}
          editable={editable}
          onClose={() => setCommentRelationship(null)}
        />
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
