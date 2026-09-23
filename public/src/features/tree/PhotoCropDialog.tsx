import { useEffect, useRef, useState } from "react";
import { Button, Modal, Notice } from "../../shared/ui";

const WIDTH = 420,
  HEIGHT = 600;
export function PhotoCropDialog({
  file,
  onCancel,
  onConfirm,
}: {
  file: File;
  onCancel: () => void;
  onConfirm: (file: File) => void;
}) {
  const canvas = useRef<HTMLCanvasElement>(null);
  const image = useRef<HTMLImageElement>();
  const drag = useRef<{ x: number; y: number; panX: number; panY: number }>();
  const [ready, setReady] = useState(false);
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  useEffect(() => {
    const url = URL.createObjectURL(file);
    const img = new Image();
    let active = true;
    img.onload = () => {
      if (active) {
        image.current = img;
        setReady(true);
      }
    };
    img.onerror = () => {
      if (active) setError("Файл не удалось прочитать как изображение");
    };
    img.src = url;
    return () => {
      active = false;
      URL.revokeObjectURL(url);
    };
  }, [file]);
  useEffect(() => {
    const img = image.current;
    const context = canvas.current?.getContext("2d");
    if (!img || !context || !ready) return;
    const scale =
      Math.max(WIDTH / img.naturalWidth, HEIGHT / img.naturalHeight) * zoom;
    const width = img.naturalWidth * scale,
      height = img.naturalHeight * scale;
    const x = Math.max(
      -(width - WIDTH) / 2,
      Math.min((width - WIDTH) / 2, pan.x),
    );
    const y = Math.max(
      -(height - HEIGHT) / 2,
      Math.min((height - HEIGHT) / 2, pan.y),
    );
    context.clearRect(0, 0, WIDTH, HEIGHT);
    context.drawImage(
      img,
      (WIDTH - width) / 2 + x,
      (HEIGHT - height) / 2 + y,
      width,
      height,
    );
  }, [ready, zoom, pan]);
  return (
    <Modal title="Настроить фотографию" onClose={onCancel} busy={busy}>
      <div className="crop-form">
        {error && <Notice>{error}</Notice>}
        <canvas
          ref={canvas}
          width={WIDTH}
          height={HEIGHT}
          aria-label="Предпросмотр фотографии"
          style={{
            width: "min(280px, 100%)",
            aspectRatio: "21/30",
            touchAction: "none",
            margin: "auto",
          }}
          onPointerDown={(event) => {
            event.currentTarget.setPointerCapture(event.pointerId);
            drag.current = {
              x: event.clientX,
              y: event.clientY,
              panX: pan.x,
              panY: pan.y,
            };
          }}
          onPointerMove={(event) => {
            if (drag.current) {
              const ratio =
                WIDTH / event.currentTarget.getBoundingClientRect().width;
              setPan({
                x: drag.current.panX + (event.clientX - drag.current.x) * ratio,
                y: drag.current.panY + (event.clientY - drag.current.y) * ratio,
              });
            }
          }}
          onPointerUp={() => {
            drag.current = undefined;
          }}
          onPointerCancel={() => {
            drag.current = undefined;
          }}
        />
        <label>
          Масштаб
          <input
            type="range"
            min="1"
            max="3"
            step="0.01"
            value={zoom}
            onChange={(event) => setZoom(Number(event.target.value))}
          />
        </label>
        <div className="modal-actions">
          <Button variant="secondary" disabled={busy} onClick={onCancel}>
            Отмена
          </Button>
          <Button
            disabled={!ready}
            loading={busy}
            onClick={() => {
              setBusy(true);
              canvas.current?.toBlob(
                (blob) => {
                  setBusy(false);
                  if (!blob) {
                    setError("Не удалось обработать фотографию");
                    return;
                  }
                  onConfirm(
                    new File([blob], file.name.replace(/\.[^.]+$/, ".jpg"), {
                      type: "image/jpeg",
                    }),
                  );
                },
                "image/jpeg",
                0.92,
              );
            }}
          >
            Применить фото
          </Button>
        </div>
      </div>
    </Modal>
  );
}
