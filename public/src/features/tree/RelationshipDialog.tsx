import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { api } from "../../api";
import type {
  Person,
  Relationship,
  RelationshipInput,
  Tree,
} from "../../types";
import { fullName } from "../../utils";
import { useTreeCache } from "../../data/queries";
import { Button, Field, Modal, Notice, Select } from "../../shared/ui";
import { relationshipProblem } from "../../graph/model";

export function RelationshipDialog({
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
  const [type, setType] = useState<RelationshipInput["type"]>("parent_child");
  const [one, setOne] = useState(selected?.id ?? "");
  const [two, setTwo] = useState("");
  const [error, setError] = useState("");
  const cache = useTreeCache(tree.id);
  const mutation = useMutation({
    mutationFn: (body: RelationshipInput) =>
      api.createRelationship(tree.id, body),
    onSuccess: async () => {
      await cache.refresh();
      onSaved();
    },
    onError: (reason) => setError(reason.message),
  });
  return (
    <Modal title="Добавить связь" onClose={onClose} busy={mutation.isPending}>
      <form
        className="modal-form"
        onSubmit={(event) => {
          event.preventDefault();
          if (mutation.isPending) return;
          const body: RelationshipInput = {
            person1_id: one,
            person2_id: two,
            type,
            ...(type === "parent_child"
              ? { direction: "parent" as const }
              : {}),
          };
          const problem = relationshipProblem(people, relationships, body);
          if (problem) return setError(problem);
          setError("");
          mutation.mutate(body);
        }}
      >
        {error && <Notice>{error}</Notice>}
        <Field label="Тип связи">
          <Select
            value={type}
            onChange={(event) =>
              setType(event.target.value as RelationshipInput["type"])
            }
          >
            <option value="parent_child">Родитель → ребёнок</option>
            <option value="spouse">Партнёры / супруги</option>
          </Select>
        </Field>
        <Field label={type === "spouse" ? "Первый человек" : "Родитель"}>
          <Select
            value={one}
            onChange={(event) => {
              setOne(event.target.value);
              if (two === event.target.value) setTwo("");
            }}
          >
            <option value="">Выберите человека</option>
            {people.map((person) => (
              <option key={person.id} value={person.id}>
                {fullName(person)}
              </option>
            ))}
          </Select>
        </Field>
        <Field label={type === "spouse" ? "Второй человек" : "Ребёнок"}>
          <Select value={two} onChange={(event) => setTwo(event.target.value)}>
            <option value="">Выберите человека</option>
            {people
              .filter((person) => person.id !== one)
              .map((person) => (
                <option key={person.id} value={person.id}>
                  {fullName(person)}
                </option>
              ))}
          </Select>
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
            Добавить связь
          </Button>
        </div>
      </form>
    </Modal>
  );
}
