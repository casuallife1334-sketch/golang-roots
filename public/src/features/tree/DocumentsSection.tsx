import { useRef, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Download, FileText, Trash2, Upload } from "lucide-react";
import { api } from "../../api";
import { documentsQuery, keys } from "../../data/queries";
import { Button, Notice } from "../../shared/ui";
import type { DocumentOwnerType } from "../../types";
import { useAuth } from "../../auth";

const maxDocumentSize = 20 * 1024 * 1024;
const acceptedTypes = [
  "application/pdf",
  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  "image/jpeg",
  "image/png",
  "image/webp",
].join(",");

export function DocumentsSection({
  treeId,
  ownerType,
  ownerId,
  editable,
}: {
  treeId: string;
  ownerType: DocumentOwnerType;
  ownerId: string;
  editable: boolean;
}) {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const input = useRef<HTMLInputElement>(null);
  const [error, setError] = useState("");
  const query = useQuery(documentsQuery(user?.id, treeId, ownerType, ownerId));
  const upload = useMutation({
    mutationFn: async (files: File[]) => {
      for (const file of files) {
        if (file.size > maxDocumentSize)
          throw new Error("Файл слишком большой. Максимальный размер — 20 МБ");
        await api.uploadDocument(treeId, ownerType, ownerId, file);
      }
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: keys.documents(user?.id, treeId, ownerType, ownerId),
      });
    },
    onError: (reason) => setError(reason instanceof Error ? reason.message : "Не удалось загрузить документ"),
  });
  const remove = useMutation({
    mutationFn: (id: string) => api.deleteDocument(treeId, id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: keys.documents(user?.id, treeId, ownerType, ownerId),
      });
    },
    onError: (reason) => setError(reason instanceof Error ? reason.message : "Не удалось удалить документ"),
  });
  const download = useMutation({
    mutationFn: (id: string) => api.downloadDocument(treeId, id),
    onSuccess: (blob, id) => {
      const document = query.data?.find((item) => item.id === id);
      if (!document) return;
      const url = URL.createObjectURL(blob);
      const link = window.document.createElement("a");
      link.href = url;
      link.download = document.file_name;
      link.click();
      URL.revokeObjectURL(url);
    },
    onError: (reason) => setError(reason instanceof Error ? reason.message : "Не удалось скачать документ"),
  });
  const documents = query.data ?? [];

  return (
    <section className="documents-section">
      <div className="documents-section-head">
        <div className="comment-editor-heading">
          <span className="comment-editor-icon">
            <FileText size={17} />
          </span>
          <span>
            <strong>Документы</strong>
            <small>Материалы, связанные с этой записью</small>
          </span>
        </div>
        {editable && (
          <>
            <input
              ref={input}
              hidden
              type="file"
              multiple
              accept={acceptedTypes}
              onChange={(event) => {
                const files = [...(event.target.files ?? [])];
                event.target.value = "";
                if (!files.length) return;
                setError("");
                upload.mutate(files);
              }}
            />
            <Button
              variant="secondary"
              loading={upload.isPending}
              onClick={() => input.current?.click()}
            >
              <Upload size={15} />
              Добавить
            </Button>
          </>
        )}
      </div>
      {error && <Notice>{error}</Notice>}
      {query.isPending ? (
        <p className="documents-empty">Загрузка документов…</p>
      ) : query.error ? (
        <Notice>{query.error.message}</Notice>
      ) : documents.length ? (
        <div className="documents-list">
          {documents.map((document) => (
            <div className="document-row" key={document.id}>
              <FileText size={17} />
              <span className="document-info">
                <strong title={document.file_name}>{document.file_name}</strong>
                <small>{formatSize(document.size_bytes)}</small>
              </span>
              <button
                className="document-action"
                aria-label={`Скачать ${document.file_name}`}
                disabled={download.isPending}
                onClick={() => download.mutate(document.id)}
              >
                <Download size={15} />
              </button>
              {editable && (
                <button
                  className="document-action danger"
                  aria-label={`Удалить ${document.file_name}`}
                  disabled={remove.isPending}
                  onClick={() => remove.mutate(document.id)}
                >
                  <Trash2 size={15} />
                </button>
              )}
            </div>
          ))}
        </div>
      ) : (
        <p className="documents-empty">Документы пока не добавлены</p>
      )}
    </section>
  );
}

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} Б`;
  if (bytes < 1024 * 1024) return `${Math.round(bytes / 1024)} КБ`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} МБ`;
}
