import { useRef, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { ImagePlus, Trash2 } from "lucide-react";
import { Button, Field, Modal, Notice, Select } from "../../shared/ui";
import { useTreeCache } from "../../data/queries";
import { dateInput, useObjectUrl, validImage } from "../../utils";
import type { Person, Tree } from "../../types";
import { personInput, PersonSave, type PersonForm } from "./personSave";
import { PhotoCropDialog } from "./PhotoCropDialog";
import { PersonPortrait } from "../../shared/PersonPortrait";
import { CommentEditor } from "./CommentEditor";

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
  const photoInput = useRef<HTMLInputElement>(null);
  const [form, setForm] = useState<PersonForm>({
    first_name: person?.first_name ?? "",
    last_name: person?.last_name ?? "",
    patronymic: person?.patronymic ?? "",
    birth_date: dateInput(person?.birth_date),
    death_date: dateInput(person?.death_date),
    gender: person?.gender ?? "",
    city: person?.metadata?.city ?? "",
    occupation: person?.metadata?.occupation ?? "",
    comment: person?.metadata?.comment ?? "",
  });
  const [file, setFile] = useState<File>();
  const [cropFile, setCropFile] = useState<File>();
  const [removePhoto, setRemovePhoto] = useState(false);
  const preview = useObjectUrl(file);
  const hasExistingPhoto = Boolean(person?.photo_url && !removePhoto);
  const hasPhoto = Boolean(file || hasExistingPhoto);
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
        wide
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
          <div className="person-form-layout">
            <fieldset
              disabled={mutation.isPending}
              className="form-fields person-main-fields"
            >
              <div className="person-name-grid">
                {field("first_name", "Имя", "text", true)}
                {field("patronymic", "Отчество")}
                {field("last_name", "Фамилия", "text", true)}
              </div>
              <div className="form-grid person-date-grid">
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
              <div className="field person-photo-field">
                <span>Фотография</span>
                <input
                  ref={photoInput}
                  hidden
                  type="file"
                  aria-label="Фотография"
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
                {hasPhoto ? (
                  <div className="person-photo-control">
                    {file && preview ? (
                      <img
                        className="photo-control-preview"
                        src={preview}
                        alt="Предпросмотр выбранной фотографии"
                      />
                    ) : (
                      person && (
                        <PersonPortrait
                          person={person}
                          treeId={tree.id}
                          className="photo-control-preview"
                        />
                      )
                    )}
                    <span className="person-photo-copy">
                      <strong>
                        {file
                          ? "Фотография подготовлена"
                          : "Фотография добавлена"}
                      </strong>
                      <small>{file?.name || "Фото профиля"}</small>
                    </span>
                    <button
                      type="button"
                      className="person-photo-remove"
                      onClick={() => {
                        setFile(undefined);
                        setRemovePhoto(Boolean(person?.photo_url));
                      }}
                    >
                      <Trash2 size={14} />
                      Удалить
                    </button>
                  </div>
                ) : (
                  <button
                    type="button"
                    className="person-photo-upload"
                    aria-label="Загрузить фотографию"
                    onClick={() => photoInput.current?.click()}
                  >
                    <span className="person-photo-upload-icon">
                      <ImagePlus size={17} />
                    </span>
                    <span>
                      <strong>Загрузить фотографию</strong>
                      <small>JPEG, PNG или WebP до 9,5 МБ</small>
                    </span>
                  </button>
                )}
                {removePhoto && (
                  <small className="photo-remove-note">
                    Фотография будет удалена после сохранения
                  </small>
                )}
              </div>
            </fieldset>
            <CommentEditor
              className="person-comment-section"
              value={form.comment}
              onChange={(comment) => setForm({ ...form, comment })}
              description="Воспоминания и важные сведения о человеке"
              placeholder="Добавьте воспоминания, биографические сведения или важные заметки о человеке…"
              disabled={mutation.isPending}
            />
          </div>
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
