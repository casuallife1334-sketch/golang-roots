import { MessageSquareText } from "lucide-react";

export function CommentEditor({
  value,
  onChange,
  title = "Комментарий",
  description,
  placeholder,
  disabled = false,
  editable = true,
  className = "",
}: {
  value: string;
  onChange?: (value: string) => void;
  title?: string;
  description?: string;
  placeholder: string;
  disabled?: boolean;
  editable?: boolean;
  className?: string;
}) {
  const length = Array.from(value).length;
  return (
    <section className={`comment-editor ${className}`}>
      <div className="comment-editor-heading">
        <span className="comment-editor-icon">
          <MessageSquareText size={17} />
        </span>
        <span>
          <strong>{title}</strong>
          {description && <small>{description}</small>}
        </span>
      </div>
      {editable ? (
        <>
          <textarea
            className="comment-editor-textarea"
            aria-label={title}
            placeholder={placeholder}
            disabled={disabled}
            value={value}
            onChange={(event) =>
              onChange?.(Array.from(event.target.value).slice(0, 5000).join(""))
            }
          />
          <span className="comment-counter" aria-live="polite">
            {length.toLocaleString("ru-RU")} / 5 000
          </span>
        </>
      ) : (
        <p className="comment-editor-readonly" tabIndex={0}>
          {value}
        </p>
      )}
    </section>
  );
}
