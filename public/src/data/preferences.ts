import { useEffect, useState } from "react";
import { z } from "zod";
import { useAuth } from "../auth";
const schema = z.object({
  portraits: z.boolean(),
  minimap: z.boolean(),
  lineWidth: z.enum(["thin", "medium"]),
});
export const defaults = {
  portraits: true,
  minimap: true,
  lineWidth: "medium" as const,
};
export type Preferences = z.infer<typeof schema>;
export function readPreferences(key: string): Preferences {
  try {
    return schema.parse({
      ...defaults,
      ...JSON.parse(localStorage.getItem(key) ?? "{}"),
    });
  } catch {
    return defaults;
  }
}
export function usePreferences() {
  const { user } = useAuth();
  const key = `roots:preferences:v1:${user?.id}`;
  const [prefs, setPrefs] = useState<Preferences>(() => readPreferences(key));
  useEffect(() => {
    const update = () => setPrefs(readPreferences(key));
    update();
    window.addEventListener("storage", update);
    window.addEventListener("roots:preferences", update);
    return () => {
      window.removeEventListener("storage", update);
      window.removeEventListener("roots:preferences", update);
    };
  }, [key]);
  return prefs;
}
