import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { FileText } from "lucide-react";
import { api } from "../../api";
import { useTreeCache } from "../../data/queries";
import { parentChild } from "../../graph/model";
import { Button, Modal, Notice } from "../../shared/ui";
import type { Person, Relationship, Tree } from "../../types";
import { fullName, metadataComment } from "../../utils";
import { CommentEditor } from "./CommentEditor";

export function RelationshipDetailsDialog({
  tree,
  relationship,
  people,
  editable,
  onClose,
}: {
  tree: Tree;
  relationship: Relationship;
  people: Person[];
  editable: boolean;
  onClose: () => void;
}) {
  const [comment, setComment] = useState(
    metadataComment(relationship.metadata),
  );
  const cache = useTreeCache(tree.id);
  const link = parentChild(relationship);
  const first = people.find(
    (person) => person.id === (link?.parentId ?? relationship.person1_id),
  );
  const second = people.find(
    (person) => person.id === (link?.childId ?? relationship.person2_id),
  );
  const description = link
    ? `${fullName(first)} → ${fullName(second)}`
    : `${fullName(first)} и ${fullName(second)}`;
  const mutation = useMutation({
    mutationFn: () => {
      const metadata = { ...relationship.metadata };
      const value = comment.trim();
      if (value) metadata.comment = value;
      else delete metadata.comment;
      return api.patchRelationship(tree.id, relationship.id, { metadata });
    },
    onSuccess: (saved) => {
      cache.saveRelationship(saved);
      onClose();
    },
  });

  return (
    <Modal
      title={
        <span className="relationship-dialog-title">
          <strong>Сведения о связи</strong>
          <small>{description}</small>
        </span>
      }
      onClose={onClose}
      busy={mutation.isPending}
      wide
    >
      {mutation.error && <Notice>{mutation.error.message}</Notice>}
      <div className="relationship-details-layout">
        <CommentEditor
          className="relationship-comment-editor"
          value={comment}
          onChange={setComment}
          description={
            editable
              ? "Сведения относятся только к этой связи"
              : "Сведения об этой связи"
          }
          placeholder="Добавьте важные сведения об отношениях между людьми…"
          disabled={mutation.isPending}
          editable={editable}
        />
        <section className="relationship-documents-panel">
          <div className="comment-editor-heading">
            <span className="comment-editor-icon">
              <FileText size={17} />
            </span>
            <span>
              <strong>Документы</strong>
              <small>Материалы, связанные с этой связью</small>
            </span>
          </div>
          <div className="relationship-documents-empty">
            Возможность добавлять документы скоро появится
          </div>
        </section>
      </div>
      <div className="modal-actions">
        <Button
          variant="secondary"
          disabled={mutation.isPending}
          onClick={onClose}
        >
          {editable ? "Отмена" : "Закрыть"}
        </Button>
        {editable && (
          <Button
            loading={mutation.isPending}
            onClick={() => mutation.mutate()}
          >
            Сохранить
          </Button>
        )}
      </div>
    </Modal>
  );
}
