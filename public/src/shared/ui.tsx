import { AlertCircle, Check, ChevronDown, LoaderCircle, X } from "lucide-react";
import type { ReactNode } from "react";
import * as Dialog from "@radix-ui/react-dialog";
import "./dialog.css";

export function Button({
  children,
  variant = "primary",
  loading,
  ...props
}: React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost" | "danger";
  loading?: boolean;
}) {
  return (
    <button
      {...props}
      type={props.type ?? "button"}
      className={`button ${variant} ${props.className ?? ""}`}
      disabled={loading || props.disabled}
    >
      {loading && <LoaderCircle className="spin" size={16} />}
      {children}
    </button>
  );
}
export function Field({
  label,
  error,
  children,
  hint,
}: {
  label: string;
  error?: string;
  children: ReactNode;
  hint?: string;
}) {
  return (
    <label className="field">
      <span>{label}</span>
      {children}
      {hint && <small className="field-hint">{hint}</small>}
      {error && <small className="field-error">{error}</small>}
    </label>
  );
}
export function Notice({
  children,
  type = "error",
}: {
  children: ReactNode;
  type?: "error" | "success";
}) {
  return (
    <div
      className={`notice ${type}`}
      role={type === "error" ? "alert" : "status"}
    >
      <AlertCircle size={16} />
      {children}
    </div>
  );
}
export function EmptyState({
  icon,
  title,
  text,
  action,
}: {
  icon?: ReactNode;
  title: string;
  text: string;
  action?: ReactNode;
}) {
  return (
    <div className="empty-state">
      {icon || <span className="empty-mark">+</span>}
      <h2>{title}</h2>
      <p>{text}</p>
      {action}
    </div>
  );
}
export function Modal({
  title,
  children,
  onClose,
  wide = false,
  busy = false,
}: {
  title: string;
  children: ReactNode;
  onClose: () => void;
  wide?: boolean;
  busy?: boolean;
}) {
  return (
    <Dialog.Root
      open
      onOpenChange={(open) => {
        if (!open && !busy) onClose();
      }}
    >
      <Dialog.Portal>
        <Dialog.Overlay className="modal-backdrop">
          <Dialog.Content
            className={`modal ${wide ? "wide" : ""}`}
            aria-describedby={undefined}
          >
            <div className="modal-head">
              <Dialog.Title>{title}</Dialog.Title>
              <button
                className="icon-button"
                onClick={onClose}
                disabled={busy}
                aria-label="Закрыть"
              >
                <X size={18} />
              </button>
            </div>
            {children}
          </Dialog.Content>
        </Dialog.Overlay>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
export function Select({
  value,
  onChange,
  children,
  ...props
}: React.SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <div className="select-wrap">
      <select {...props} value={value} onChange={onChange}>
        {children}
      </select>
      <ChevronDown size={15} />
    </div>
  );
}
export function SavedBadge() {
  return (
    <span className="saved-badge">
      <Check size={13} />
      Сохранено
    </span>
  );
}
