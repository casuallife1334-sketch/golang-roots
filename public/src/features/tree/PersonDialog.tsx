import { useRef, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { Button, Field, Modal, Notice, Select } from "../../shared/ui";
import { useTreeCache } from "../../data/queries";
import { dateInput, validImage } from "../../utils";
import type { Person, Tree } from "../../types";
import { personInput, PersonSave, type PersonForm } from "./personSave";
import { PhotoCropDialog } from "./PhotoCropDialog";

export function PersonDialog({
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
  const cache = useTreeCache(tree.id);
  const save = useRef(new PersonSave(person));
  const [form, setForm] = useState<PersonForm>({
    first_name: person?.first_name ?? "",
    last_name: person?.last_name ?? "",
    patronymic: person?.metadata?.patronymic ?? "",
    birth_date: dateInput(person?.birth_date),
    death_date: dateInput(person?.death_date),
    gender: person?.gender ?? "",
    city: person?.metadata?.city ?? "",
    occupation: person?.metadata?.occupation ?? "",
  });
  const [file, setFile] = useState<File>();
  const [cropFile, setCropFile] = useState<File>();
  const [removePhoto, setRemovePhoto] = useState(false);
  const [error, setError] = useState("");
  const mutation = useMutation({
    mutationFn: () =>
      save.current.run(
        tree.id,
        personInput(form, person),
        file,
        removePhoto,
        cache.savePerson,
      ),
    onSuccess: (saved) => {
      void cache.refreshPhoto(saved.id);
      void cache.refresh();
      onSaved(saved);
    },
    onError: (reason) =>
      setError(
        reason instanceof Error
          ? reason.message
          : "Не удалось сохранить человека",
      ),
  });
  const field = (
    key: keyof PersonForm,
    label: string,
    type = "text",
    required = false,
  ) => (
    <Field label={label}>
      <input
        type={type}
        required={required}
        maxLength={200}
        value={form[key]}
        onChange={(event) => setForm({ ...form, [key]: event.target.value })}
      />
    </Field>
  );
  return (
    <>
      <Modal
        title={person ? "Редактировать человека" : "Добавить человека"}
        onClose={onClose}
        busy={mutation.isPending}
      >
        <form
          className="modal-form"
          onSubmit={(event) => {
            event.preventDefault();
            if (mutation.isPending) return;
            setError("");
            if (!form.first_name.trim() || !form.last_name.trim())
              return setError("Укажите имя и фамилию");
            if (
              form.birth_date &&
              form.death_date &&
              form.birth_date > form.death_date
            )
              return setError("Дата рождения не может быть позже даты смерти");
            mutation.mutate();
          }}
        >
          {error && <Notice>{error}</Notice>}
          <fieldset disabled={mutation.isPending} className="form-fields">
            <div className="form-grid">
              {field("first_name", "Имя", "text", true)}
              {field("patronymic", "Отчество")}
              {field("last_name", "Фамилия", "text", true)}
              {field("birth_date", "Дата рождения", "date")}
              {field("death_date", "Дата смерти", "date")}
            </div>
            <Field label="Пол">
              <Select
                value={form.gender}
                onChange={(event) =>
                  setForm({ ...form, gender: event.target.value })
                }
              >
                <option value="">Не указан</option>
                <option value="female">Женский</option>
                <option value="male">Мужской</option>
                <option value="other">Другой</option>
              </Select>
            </Field>
            <div className="form-grid">
              {field("city", "Место / город")}
              {field("occupation", "Род деятельности")}
            </div>
            <Field label="Фотография">
              <input
                type="file"
                accept="image/jpeg,image/png,image/webp"
                onChange={(event) => {
                  const next = event.target.files?.[0];
                  event.target.value = "";
                  if (!next) return;
                  if (!validImage(next))
                    return setError(
                      "Поддерживаются JPEG, PNG и WebP до 9,5 МБ",
                    );
                  setCropFile(next);
                  setRemovePhoto(false);
                }}
              />
            </Field>
            {file && <p>{file.name}</p>}
            {(file || person?.photo_url) && (
              <Button
                variant="ghost"
                onClick={() => {
                  setFile(undefined);
                  setRemovePhoto(true);
                }}
              >
                Удалить фотографию
              </Button>
            )}
            {removePhoto && <p>Фотография будет удалена</p>}
          </fieldset>
          <div className="modal-actions">
            <Button
              variant="secondary"
              disabled={mutation.isPending}
              onClick={onClose}
            >
              Закрыть
            </Button>
            <Button type="submit" loading={mutation.isPending}>
              {person || save.current.saved ? "Сохранить" : "Добавить человека"}
            </Button>
          </div>
        </form>
      </Modal>
      {cropFile && (
        <PhotoCropDialog
          file={cropFile}
          onCancel={() => setCropFile(undefined)}
          onConfirm={(value) => {
            setFile(value);
            setCropFile(undefined);
          }}
        />
      )}
    </>
  );
}
