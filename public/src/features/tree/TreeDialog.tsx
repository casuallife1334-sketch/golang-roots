import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "../../api";
import { useAuth } from "../../auth";
import { keys } from "../../data/queries";
import { Button, Field, Modal, Notice } from "../../shared/ui";
export function TreeDialog({
  onClose,
  onSaved,
}: {
  onClose: () => void;
  onSaved: (id: string) => void;
}) {
  const [name, setName] = useState("");
  const { user } = useAuth();
  const query = useQueryClient();
  const mutation = useMutation({
    mutationFn: () => api.createTree(name.trim()),
    onSuccess: async (tree) => {
      await query.invalidateQueries({ queryKey: keys.trees(user?.id) });
      onSaved(tree.id);
    },
  });
  return (
    <Modal title="Новое дерево" onClose={onClose} busy={mutation.isPending}>
      <form
        className="modal-form"
        onSubmit={(event) => {
          event.preventDefault();
          if (name.trim() && !mutation.isPending) mutation.mutate();
        }}
      >
        {mutation.error && <Notice>{mutation.error.message}</Notice>}
        <Field label="Название дерева">
          <input
            autoFocus
            required
            maxLength={200}
            value={name}
            onChange={(event) => setName(event.target.value)}
          />
        </Field>
        <div className="modal-actions">
          <Button
            variant="secondary"
            disabled={mutation.isPending}
            onClick={onClose}
          >
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
