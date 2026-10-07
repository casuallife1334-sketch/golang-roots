import type { z } from "zod";

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public code?: string,
  ) {
    super(message);
  }
}
export const TOKEN_KEY = "roots:access-token";
let unauthorized: (() => void) | undefined;
export function setUnauthorizedHandler(handler: () => void) {
  unauthorized = handler;
  return () => {
    if (unauthorized === handler) unauthorized = undefined;
  };
}
let sessionRequests = new AbortController();
export function cancelSessionRequests() {
  sessionRequests.abort();
  sessionRequests = new AbortController();
}
async function validatedImageBlob(response: Response) {
  const buffer = await response.arrayBuffer();
  const bytes = new Uint8Array(buffer, 0, Math.min(12, buffer.byteLength));
  const type =
    bytes[0] === 0xff && bytes[1] === 0xd8 && bytes[2] === 0xff
      ? "image/jpeg"
      : bytes.length >= 8 &&
          [0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a].every(
            (value, index) => bytes[index] === value,
          )
        ? "image/png"
        : bytes.length >= 12 &&
            String.fromCharCode(...bytes.slice(0, 4)) === "RIFF" &&
            String.fromCharCode(...bytes.slice(8, 12)) === "WEBP"
          ? "image/webp"
          : null;
  if (!type)
    throw new ApiError(
      502,
      "Сервер вернул некорректную фотографию",
      "invalid_response",
    );
  return new Blob([buffer], { type });
}
export async function request<T>(
  path: string,
  schema: z.ZodType<T, z.ZodTypeDef, unknown> | "blob" | "file" | "void",
  init: RequestInit = {},
  anonymous = false,
): Promise<T> {
  const token = anonymous ? null : sessionStorage.getItem(TOKEN_KEY);
  const headers = new Headers(init.headers);
  if (init.body && !(init.body instanceof FormData))
    headers.set("Content-Type", "application/json");
  if (token) headers.set("Authorization", `Bearer ${token}`);
  const controller = new AbortController();
  const abort = () => controller.abort();
  const signals = [init.signal, sessionRequests.signal].filter(
    (s): s is AbortSignal => Boolean(s),
  );
  signals.forEach((signal) => {
    signal.addEventListener("abort", abort, { once: true });
    if (signal.aborted) abort();
  });
  let timedOut = false;
  const timer = setTimeout(() => {
    timedOut = true;
    controller.abort();
  }, 30_000);
  try {
    const response = await fetch(`/api/v1${path}`, {
      ...init,
      headers,
      signal: controller.signal,
    });
    if (!response.ok) {
      if (
        response.status === 401 &&
        !anonymous &&
        token === sessionStorage.getItem(TOKEN_KEY)
      )
        unauthorized?.();
      const body = await response.json().catch(() => null);
      const serverMessage =
        typeof body?.message === "string"
          ? body.message
          : typeof body?.error === "string"
            ? body.error
            : "";
      const messages: Record<number, string> = {
        400: serverMessage || "Проверьте данные формы",
        401: anonymous
          ? "Неверный email или пароль"
          : "Сессия истекла. Войдите снова",
        403: "У вас нет доступа к этому действию",
        404: "Запись не найдена",
        409: "Такая запись уже существует",
        413: "Фотография слишком большая. Максимальный размер — 9,5 МБ",
      };
      throw new ApiError(
        response.status,
        messages[response.status] || "Сервис временно недоступен",
      );
    }
    if (schema === "void") return undefined as T;
    if (schema === "blob") {
      return (await validatedImageBlob(response)) as T;
    }
    if (schema === "file") return (await response.blob()) as T;
    const parsed = schema.safeParse(
      await response.json().catch(() => undefined),
    );
    if (!parsed.success)
      throw new ApiError(
        502,
        "Сервер вернул данные в неожиданном формате",
        "invalid_response",
      );
    return parsed.data;
  } catch (error) {
    if (error instanceof ApiError) throw error;
    if (timedOut)
      throw new ApiError(
        0,
        "Сервер не ответил вовремя. Повторите запрос",
        "timeout",
      );
    if (controller.signal.aborted)
      throw new DOMException("Запрос отменён", "AbortError");
    throw new ApiError(0, "Не удалось подключиться к серверу");
  } finally {
    clearTimeout(timer);
    signals.forEach((signal) => signal.removeEventListener("abort", abort));
  }
}
