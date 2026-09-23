import { useEffect, useRef, useState, type CSSProperties } from "react";
import { Crop, Minus, Move, Plus, RotateCcw } from "lucide-react";
import { Button, Modal, Notice } from "../../shared/ui";

const OUTPUT_WIDTH = 420;
const OUTPUT_HEIGHT = 600;
const STAGE_SIZE = 560;
const CROP_WIDTH = 280;
const CROP_HEIGHT = 400;
const CROP_X = (STAGE_SIZE - CROP_WIDTH) / 2;
const CROP_Y = (STAGE_SIZE - CROP_HEIGHT) / 2;

type Point = { x: number; y: number };

function imageScale(image: HTMLImageElement, zoom: number) {
  return (
    Math.max(
      CROP_WIDTH / image.naturalWidth,
      CROP_HEIGHT / image.naturalHeight,
    ) * zoom
  );
}

function boundedPan(image: HTMLImageElement, zoom: number, pan: Point): Point {
  const scale = imageScale(image, zoom);
  const limitX = Math.max(0, (image.naturalWidth * scale - CROP_WIDTH) / 2);
  const limitY = Math.max(0, (image.naturalHeight * scale - CROP_HEIGHT) / 2);
  return {
    x: Math.max(-limitX, Math.min(limitX, pan.x)),
    y: Math.max(-limitY, Math.min(limitY, pan.y)),
  };
}

function drawCropGuides(context: CanvasRenderingContext2D) {
  context.fillStyle = "rgb(17 27 45 / 58%)";
  context.fillRect(0, 0, STAGE_SIZE, CROP_Y);
  context.fillRect(0, CROP_Y + CROP_HEIGHT, STAGE_SIZE, CROP_Y);
  context.fillRect(0, CROP_Y, CROP_X, CROP_HEIGHT);
  context.fillRect(CROP_X + CROP_WIDTH, CROP_Y, CROP_X, CROP_HEIGHT);

  context.strokeStyle = "rgb(255 255 255 / 38%)";
  context.lineWidth = 1;
  for (const part of [1, 2]) {
    const x = CROP_X + (CROP_WIDTH / 3) * part;
    const y = CROP_Y + (CROP_HEIGHT / 3) * part;
    context.beginPath();
    context.moveTo(x, CROP_Y);
    context.lineTo(x, CROP_Y + CROP_HEIGHT);
    context.moveTo(CROP_X, y);
    context.lineTo(CROP_X + CROP_WIDTH, y);
    context.stroke();
  }

  context.strokeStyle = "#ffffff";
  context.lineWidth = 2;
  context.strokeRect(CROP_X, CROP_Y, CROP_WIDTH, CROP_HEIGHT);
  context.strokeStyle = "#3478ff";
  context.lineWidth = 6;
  const corner = 25;
  const corners: [number, number, number, number][] = [
    [CROP_X, CROP_Y + corner, CROP_X, CROP_Y],
    [CROP_X, CROP_Y, CROP_X + corner, CROP_Y],
    [CROP_X + CROP_WIDTH - corner, CROP_Y, CROP_X + CROP_WIDTH, CROP_Y],
    [CROP_X + CROP_WIDTH, CROP_Y, CROP_X + CROP_WIDTH, CROP_Y + corner],
    [
      CROP_X + CROP_WIDTH,
      CROP_Y + CROP_HEIGHT - corner,
      CROP_X + CROP_WIDTH,
      CROP_Y + CROP_HEIGHT,
    ],
    [
      CROP_X + CROP_WIDTH,
      CROP_Y + CROP_HEIGHT,
      CROP_X + CROP_WIDTH - corner,
      CROP_Y + CROP_HEIGHT,
    ],
    [CROP_X + corner, CROP_Y + CROP_HEIGHT, CROP_X, CROP_Y + CROP_HEIGHT],
    [CROP_X, CROP_Y + CROP_HEIGHT, CROP_X, CROP_Y + CROP_HEIGHT - corner],
  ];
  for (const [fromX, fromY, toX, toY] of corners) {
    context.beginPath();
    context.moveTo(fromX, fromY);
    context.lineTo(toX, toY);
    context.stroke();
  }
}

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
  const drag = useRef<{ x: number; y: number; pan: Point }>();
  const [ready, setReady] = useState(false);
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState<Point>({ x: 0, y: 0 });
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    const url = URL.createObjectURL(file);
    const next = new Image();
    let active = true;
    next.onload = () => {
      if (!active) return;
      image.current = next;
      setPan({ x: 0, y: 0 });
      setZoom(1);
      setReady(true);
    };
    next.onerror = () => {
      if (active) setError("Файл не удалось прочитать как изображение");
    };
    next.src = url;
    return () => {
      active = false;
      URL.revokeObjectURL(url);
    };
  }, [file]);

  useEffect(() => {
    const source = image.current;
    const context = canvas.current?.getContext("2d");
    if (!source || !context || !ready) return;
    const scale = imageScale(source, zoom);
    const currentPan = boundedPan(source, zoom, pan);
    const width = source.naturalWidth * scale;
    const height = source.naturalHeight * scale;
    context.clearRect(0, 0, STAGE_SIZE, STAGE_SIZE);
    context.fillStyle = "#d9d7d2";
    context.fillRect(0, 0, STAGE_SIZE, STAGE_SIZE);
    context.drawImage(
      source,
      (STAGE_SIZE - width) / 2 + currentPan.x,
      (STAGE_SIZE - height) / 2 + currentPan.y,
      width,
      height,
    );
    drawCropGuides(context);
  }, [ready, zoom, pan]);

  const changeZoom = (value: number) => {
    const next = Math.max(1, Math.min(3, value));
    setZoom(next);
    if (image.current)
      setPan((current) => boundedPan(image.current!, next, current));
  };

  const exportPhoto = () => {
    const source = image.current;
    if (!source) return;
    setBusy(true);
    setError("");
    const output = document.createElement("canvas");
    output.width = OUTPUT_WIDTH;
    output.height = OUTPUT_HEIGHT;
    const context = output.getContext("2d");
    if (!context) {
      setBusy(false);
      setError("Не удалось обработать фотографию");
      return;
    }
    const ratio = OUTPUT_WIDTH / CROP_WIDTH;
    const scale = imageScale(source, zoom) * ratio;
    const currentPan = boundedPan(source, zoom, pan);
    const width = source.naturalWidth * scale;
    const height = source.naturalHeight * scale;
    context.drawImage(
      source,
      (OUTPUT_WIDTH - width) / 2 + currentPan.x * ratio,
      (OUTPUT_HEIGHT - height) / 2 + currentPan.y * ratio,
      width,
      height,
    );
    output.toBlob(
      (blob) => {
        setBusy(false);
        if (!blob) {
          setError("Не удалось обработать фотографию");
          return;
        }
        const name = file.name.replace(/\.[^.]+$/, "") || "photo";
        onConfirm(new File([blob], `${name}.jpg`, { type: "image/jpeg" }));
      },
      "image/jpeg",
      0.92,
    );
  };

  const rangeStyle = {
    "--range-progress": `${((zoom - 1) / 2) * 100}%`,
  } as CSSProperties;

  return (
    <Modal title="Настроить фотографию" onClose={onCancel} busy={busy}>
      <div className="crop-form">
        {error && <Notice>{error}</Notice>}
        <div className="crop-preview">
          <canvas
            ref={canvas}
            className="crop-canvas"
            width={STAGE_SIZE}
            height={STAGE_SIZE}
            aria-label="Область обрезки фотографии"
            onPointerDown={(event) => {
              event.currentTarget.setPointerCapture(event.pointerId);
              drag.current = {
                x: event.clientX,
                y: event.clientY,
                pan,
              };
            }}
            onPointerMove={(event) => {
              if (!drag.current || !image.current) return;
              const ratio =
                STAGE_SIZE / event.currentTarget.getBoundingClientRect().width;
              setPan(
                boundedPan(image.current, zoom, {
                  x:
                    drag.current.pan.x +
                    (event.clientX - drag.current.x) * ratio,
                  y:
                    drag.current.pan.y +
                    (event.clientY - drag.current.y) * ratio,
                }),
              );
            }}
            onPointerUp={() => {
              drag.current = undefined;
            }}
            onPointerCancel={() => {
              drag.current = undefined;
            }}
            onWheel={(event) => {
              event.preventDefault();
              changeZoom(zoom + (event.deltaY > 0 ? -0.08 : 0.08));
            }}
          />
          <span className="crop-ratio">
            <Crop size={13} /> 7:10
          </span>
        </div>
        <p className="crop-help">
          <Move size={15} /> Перетаскивайте фотографию. Затемнённые края не
          попадут в итоговый кадр.
        </p>
        <div className="crop-controls">
          <div className="crop-control-head">
            <span>Масштаб</span>
            <strong>{Math.round(zoom * 100)}%</strong>
          </div>
          <div className="crop-range-row">
            <button
              type="button"
              aria-label="Уменьшить масштаб"
              onClick={() => changeZoom(zoom - 0.1)}
              disabled={zoom <= 1}
            >
              <Minus size={15} />
            </button>
            <input
              type="range"
              min="1"
              max="3"
              step="0.01"
              value={zoom}
              style={rangeStyle}
              aria-label="Масштаб фотографии"
              onChange={(event) => changeZoom(Number(event.target.value))}
            />
            <button
              type="button"
              aria-label="Увеличить масштаб"
              onClick={() => changeZoom(zoom + 0.1)}
              disabled={zoom >= 3}
            >
              <Plus size={15} />
            </button>
          </div>
          <button
            type="button"
            className="crop-reset"
            onClick={() => {
              setZoom(1);
              setPan({ x: 0, y: 0 });
            }}
            disabled={zoom === 1 && pan.x === 0 && pan.y === 0}
          >
            <RotateCcw size={14} /> Сбросить положение
          </button>
        </div>
        <div className="modal-actions">
          <Button variant="secondary" disabled={busy} onClick={onCancel}>
            Отмена
          </Button>
          <Button disabled={!ready} loading={busy} onClick={exportPhoto}>
            Применить фото
          </Button>
        </div>
      </div>
    </Modal>
  );
}
