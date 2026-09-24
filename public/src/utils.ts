import { useEffect, useState } from "react";
import type { Person } from "./types";
export const initials = (person: Pick<Person, "first_name" | "last_name">) =>
  `${person.first_name[0] || ""}${person.last_name[0] || ""}`.toUpperCase();
export const formatDate = (date?: string | null) =>
  date
    ? new Intl.DateTimeFormat("ru-RU", {
        day: "numeric",
        month: "long",
        year: "numeric",
        timeZone: "UTC",
      }).format(new Date(date))
    : "Не указано";
export const year = (date?: string | null) =>
  date ? new Date(date).getUTCFullYear().toString() : "";
export const formatYears = (person: Person) => {
  const birth = year(person.birth_date);
  const death = year(person.death_date);
  return birth || death
    ? `${birth || "?"} — ${death || "..."}`
    : "Даты не указаны";
};
export const fullName = (person?: Person) =>
  person
    ? [person.first_name, person.patronymic, person.last_name]
        .filter(Boolean)
        .join(" ")
    : "Человек";
export const dateInput = (value?: string | null) =>
  value ? value.slice(0, 10) : "";
export const validImage = (file: File) =>
  ["image/jpeg", "image/png", "image/webp"].includes(file.type) &&
  file.size <= 9_500_000;
export function useObjectUrl(blob?: Blob) {
  const [url, setUrl] = useState<string>();
  useEffect(() => {
    if (!blob) {
      setUrl(undefined);
      return;
    }
    const next = URL.createObjectURL(blob);
    setUrl(next);
    return () => URL.revokeObjectURL(next);
  }, [blob]);
  return url;
}
