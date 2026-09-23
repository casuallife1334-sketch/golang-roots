import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  GitBranch,
  LogOut,
  Monitor,
  Palette,
  Trash2,
  UserRound,
} from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { api, ApiError } from "../api";
import { useAuth } from "../auth";
import { Button, Field, Modal, Notice, SavedBadge } from "../shared/ui";
import { formatDate } from "../utils";
import { keys, treeQuery } from "../data/queries";
import { readPreferences } from "../data/preferences";

interface Preferences {
  compact: boolean;
  portraits: boolean;
  minimap: boolean;
  lineWidth: "thin" | "medium";
}
const defaults: Preferences = {
  compact: false,
  portraits: true,
  minimap: true,
  lineWidth: "medium",
};
export function SettingsPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const query = useQueryClient();
  const trees = useQuery(treeQuery(user?.id));
  const [prefs, setPrefs] = useState(defaults);
  const [draft, setDraft] = useState(defaults);
  const [treeName, setTreeName] = useState("");
  const [modal, setModal] = useState(false);
  const [message, setMessage] = useState("");
  const key = `roots:preferences:v1:${user?.id}`;
  useEffect(() => {
    const next = readPreferences(key);
    setPrefs(next);
    setDraft(next);
  }, [key]);
  const [currentId, setCurrentId] = useState(() =>
    sessionStorage.getItem(`roots:last-tree:${user?.id}`),
  );
  const current =
    trees.data?.find((tree) => tree.id === currentId) ?? trees.data?.[0];
  useEffect(() => {
    if (current) setTreeName(current.name);
  }, [current?.id]);
  const save = () => {
    try {
      localStorage.setItem(key, JSON.stringify(draft));
    } catch {
      setMessage("Не удалось сохранить настройки браузера");
      return;
    }
    window.dispatchEvent(new Event("roots:preferences"));
    setPrefs(draft);
    setMessage("Настройки отображения сохранены");
    setTimeout(() => setMessage(""), 2400);
  };
  const rename = useMutation({
    mutationFn: () => api.patchTree(current!.id, treeName.trim()),
    onSuccess: () => {
      query.invalidateQueries({ queryKey: keys.trees(user?.id) });
      setMessage("Название дерева сохранено");
    },
  });
  const remove = useMutation({
    mutationFn: () => api.deleteTree(current!.id),
    onSuccess: () => {
      setModal(false);
      query.removeQueries({ queryKey: keys.people(user?.id, current!.id) });
      query.removeQueries({
        queryKey: keys.relationships(user?.id, current!.id),
      });
      query.removeQueries({ queryKey: ["photo", user?.id, current!.id] });
      query.invalidateQueries({ queryKey: keys.trees(user?.id) });
      navigate("/trees");
    },
  });
  return (
    <div className="content-page settings-page">
      <header className="topbar">
        <div className="topbar-title">Настройки</div>
        <div className="settings-actions">
          <Button variant="secondary" onClick={() => setDraft(prefs)}>
            Сбросить
          </Button>
          <Button onClick={save}>Сохранить изменения</Button>
        </div>
      </header>
      <section className="page-heading">
        <div>
          <h1>Настройки</h1>
          <p>Рабочая область и данные аккаунта</p>
        </div>
        {message && <p role="status">{message}</p>}
      </section>
      <div className="settings-grid">
        <section className="settings-card">
          <div className="card-title">
            <UserRound size={20} />
            <div>
              <h2>Аккаунт</h2>
              <p>Данные текущего пользователя</p>
            </div>
          </div>
          <div className="readonly-row">
            <span>Email</span>
            <strong>{user?.email}</strong>
          </div>
          <div className="readonly-row">
            <span>Аккаунт создан</span>
            <strong>{user ? formatDate(user.created_at) : "—"}</strong>
          </div>
          <button className="logout-row" onClick={logout}>
            <LogOut size={16} />
            Выйти из аккаунта
          </button>
        </section>
        <section className="settings-card">
          <div className="card-title">
            <Palette size={20} />
            <div>
              <h2>Внешний вид дерева</h2>
              <p>Настройки сохраняются в этом браузере</p>
            </div>
          </div>
          <Toggle
            label="Компактные карточки"
            description="Уменьшить расстояние между людьми"
            value={draft.compact}
            onChange={(value) => setDraft({ ...draft, compact: value })}
          />
          <Toggle
            label="Показывать портреты"
            description="Использовать загруженные фотографии"
            value={draft.portraits}
            onChange={(value) => setDraft({ ...draft, portraits: value })}
          />
          <Toggle
            label="Показывать мини-карту"
            description="Оставить обзор холста внизу"
            value={draft.minimap}
            onChange={(value) => setDraft({ ...draft, minimap: value })}
          />
        </section>
        <section className="settings-card tree-settings">
          {trees.error && <Notice>{trees.error.message}</Notice>}
          {rename.error && <Notice>{rename.error.message}</Notice>}
          <Field label="Выбрать дерево">
            <select
              value={current?.id ?? ""}
              onChange={(event) => setCurrentId(event.target.value)}
            >
              {trees.data?.map((tree) => (
                <option key={tree.id} value={tree.id}>
                  {tree.name}
                </option>
              ))}
            </select>
          </Field>
          <div className="card-title">
            <GitBranch size={20} />
            <div>
              <h2>Текущее дерево</h2>
              <p>Управление доступными деревьями</p>
            </div>
          </div>
          {current ? (
            <>
              <Field label="Название дерева">
                <div className="inline-field">
                  <input
                    value={treeName}
                    onChange={(e) => setTreeName(e.target.value)}
                  />
                  <Button
                    variant="secondary"
                    disabled={
                      current.owner_id !== user?.id ||
                      !treeName.trim() ||
                      treeName === current.name
                    }
                    loading={rename.isPending}
                    onClick={() => rename.mutate()}
                  >
                    Сохранить
                  </Button>
                </div>
              </Field>
              <div className="tree-summary">
                <span>Деревьев доступно</span>
                <strong>{trees.data?.length || 0}</strong>
              </div>
              <button
                className="danger-row"
                disabled={current.owner_id !== user?.id}
                onClick={() => setModal(true)}
              >
                <Trash2 size={16} />
                Удалить дерево
              </button>
            </>
          ) : (
            <p className="muted">
              Деревьев пока нет. Создать первое можно на странице древа.
            </p>
          )}
        </section>
        <section className="settings-card">
          <div className="card-title">
            <Monitor size={20} />
            <div>
              <h2>О Roots</h2>
              <p>Семейная история без лишнего шума</p>
            </div>
          </div>
          <p className="settings-copy">
            Данные людей, связей и деревьев хранятся через API. Настройки
            отображения сохраняются локально и не содержат личных данных.
          </p>
        </section>
      </div>
      {modal && (
        <Modal
          title="Удалить дерево?"
          onClose={() => setModal(false)}
          busy={remove.isPending}
        >
          <div className="confirm-content">
            {remove.error && <Notice>{remove.error.message}</Notice>}
            <Notice>
              Будут удалены дерево, люди и связи внутри него. Это действие
              нельзя отменить.
            </Notice>
            <p>Удалить «{current?.name}»?</p>
            <div className="modal-actions">
              <Button variant="secondary" onClick={() => setModal(false)}>
                Отмена
              </Button>
              <Button
                variant="primary"
                loading={remove.isPending}
                onClick={() => remove.mutate()}
              >
                Удалить дерево
              </Button>
            </div>
          </div>
        </Modal>
      )}
    </div>
  );
}
function Toggle({
  label,
  description,
  value,
  onChange,
}: {
  label: string;
  description: string;
  value: boolean;
  onChange: (value: boolean) => void;
}) {
  return (
    <label className="toggle-row">
      <span>
        <strong>{label}</strong>
        <small>{description}</small>
      </span>
      <input
        type="checkbox"
        checked={value}
        onChange={(e) => onChange(e.target.checked)}
      />
      <i />
    </label>
  );
}
