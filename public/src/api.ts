import type {
  Person,
  PersonInput,
  Relationship,
  RelationshipInput,
  Tree,
  User,
} from "./types";

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public code?: string,
  ) {
    super(message);
  }
}

let onUnauthorized: (() => void) | undefined;
export const setUnauthorizedHandler = (handler: () => void) => {
  onUnauthorized = handler;
};
const token = () => sessionStorage.getItem("roots:access-token");

async function request<T>(
  path: string,
  init: RequestInit = {},
  options?: { skipAuthRedirect?: boolean; loginAttempt?: boolean },
): Promise<T> {
  const headers = new Headers(init.headers);
  if (!(init.body instanceof FormData))
    headers.set("Content-Type", "application/json");
  const accessToken = token();
  if (accessToken) headers.set("Authorization", `Bearer ${accessToken}`);
  let response: Response;
  try {
    response = await fetch(`/api/v1${path}`, { ...init, headers });
  } catch {
    throw new ApiError(0, "Не удалось подключиться к серверу");
  }
  if (response.status === 204) return undefined as T;
  const contentType = response.headers.get("content-type") || "";
  const body = contentType.includes("json")
    ? await response.json().catch(() => null)
    : await response.text();
  if (!response.ok) {
    if (response.status === 401 && !options?.skipAuthRedirect)
      onUnauthorized?.();
    const message =
      typeof body === "object" && body
        ? body.message || body.error
        : "Запрос не выполнен";
    throw new ApiError(
      response.status,
      friendlyError(response.status, message, options?.loginAttempt),
    );
  }
  return body as T;
}

function friendlyError(status: number, message: string, loginAttempt = false) {
  if (status === 400) return message || "Проверьте данные формы";
  if (status === 401)
    return loginAttempt
      ? "Неверный email или пароль"
      : "Сессия истекла. Войдите снова";
  if (status === 403) return "У вас нет доступа к этому действию";
  if (status === 404) return "Запись не найдена";
  if (status === 409) return "Такая запись уже существует";
  if (status === 413)
    return "Фотография слишком большая. Максимальный размер — 9,5 МБ";
  return status >= 500
    ? "Сервис временно недоступен"
    : message || "Произошла ошибка";
}

export const api = {
  register: (body: { email: string; password: string }) =>
    request<User>(
      "/auth/register",
      { method: "POST", body: JSON.stringify(body) },
      { skipAuthRedirect: true },
    ),
  login: (body: { email: string; password: string }) =>
    request<{ access_token: string; token_type: string; expires_in: number }>(
      "/auth/login",
      { method: "POST", body: JSON.stringify(body) },
      { skipAuthRedirect: true, loginAttempt: true },
    ),
  me: () => request<User>("/users/me"),
  trees: () => request<Tree[]>("/trees"),
  tree: (id: string) => request<Tree>(`/trees/${id}`),
  createTree: (name: string) =>
    request<Tree>("/trees", { method: "POST", body: JSON.stringify({ name }) }),
  patchTree: (id: string, name: string) =>
    request<Tree>(`/trees/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ name }),
    }),
  deleteTree: (id: string) =>
    request<void>(`/trees/${id}`, { method: "DELETE" }),
  persons: (treeId: string) => request<Person[]>(`/trees/${treeId}/persons`),
  person: (treeId: string, id: string) =>
    request<Person>(`/trees/${treeId}/persons/${id}`),
  createPerson: (treeId: string, body: PersonInput) =>
    request<Person>(`/trees/${treeId}/persons`, {
      method: "POST",
      body: JSON.stringify(body),
    }),
  patchPerson: (treeId: string, id: string, body: Partial<PersonInput>) =>
    request<Person>(`/trees/${treeId}/persons/${id}`, {
      method: "PATCH",
      body: JSON.stringify(body),
    }),
  deletePerson: (treeId: string, id: string) =>
    request<void>(`/trees/${treeId}/persons/${id}`, { method: "DELETE" }),
  relationships: (treeId: string) =>
    request<Relationship[]>(`/trees/${treeId}/relationships`),
  createRelationship: (treeId: string, body: RelationshipInput) =>
    request<Relationship>(`/trees/${treeId}/relationships`, {
      method: "POST",
      body: JSON.stringify(body),
    }),
  deleteRelationship: (treeId: string, id: string) =>
    request<void>(`/trees/${treeId}/relationships/${id}`, { method: "DELETE" }),
  photo: async (treeId: string, personId: string) => {
    const accessToken = token();
    const response = await fetch(
      `/api/v1/trees/${treeId}/persons/${personId}/photo`,
      {
        headers: accessToken ? { Authorization: `Bearer ${accessToken}` } : {},
      },
    );
    if (!response.ok) {
      if (response.status === 401) onUnauthorized?.();
      throw new ApiError(response.status, "Не удалось загрузить фотографию");
    }
    return response.blob();
  },
  uploadPhoto: (treeId: string, personId: string, file: File) => {
    const body = new FormData();
    body.append("file", file);
    return request<Person>(`/trees/${treeId}/persons/${personId}/photo`, {
      method: "POST",
      body,
    });
  },
  deletePhoto: (treeId: string, personId: string) =>
    request<void>(`/trees/${treeId}/persons/${personId}/photo`, {
      method: "DELETE",
    }),
};
