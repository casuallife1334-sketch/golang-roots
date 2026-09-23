import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { useNavigate } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { api, ApiError, setUnauthorizedHandler } from "./api";
import { cancelSessionRequests, TOKEN_KEY } from "./data/http";
import type { User } from "./types";

interface AuthContextValue {
  user: User | null;
  loading: boolean;
  error: string;
  retry: () => void;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
}
const AuthContext = createContext<AuthContextValue | null>(null);
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [attempt, setAttempt] = useState(0);
  const epoch = useRef(0);
  const navigate = useNavigate();
  const navigateRef = useRef(navigate);
  navigateRef.current = navigate;
  const query = useQueryClient();
  const clear = useCallback(() => {
    epoch.current++;
    cancelSessionRequests();
    query.clear();
    sessionStorage.removeItem(TOKEN_KEY);
    setUser(null);
    setError("");
    setLoading(false);
    navigateRef.current("/login", { replace: true });
  }, [query]);
  useEffect(() => setUnauthorizedHandler(clear), [clear]);
  useEffect(() => {
    const controller = new AbortController();
    const version = ++epoch.current;
    if (!sessionStorage.getItem(TOKEN_KEY)) {
      setLoading(false);
      return;
    }
    setLoading(true);
    setError("");
    api
      .me(controller.signal)
      .then((value) => {
        if (version === epoch.current) setUser(value);
      })
      .catch((reason) => {
        if (controller.signal.aborted || version !== epoch.current) return;
        if (reason instanceof ApiError && reason.status === 401) clear();
        else
          setError(
            "Не удалось проверить сессию. Проверьте соединение и повторите попытку.",
          );
      })
      .finally(() => {
        if (version === epoch.current) setLoading(false);
      });
    return () => controller.abort();
  }, [attempt, clear]);
  const login = async (email: string, password: string) => {
    const version = ++epoch.current;
    cancelSessionRequests();
    query.clear();
    const result = await api.login({ email, password });
    if (version !== epoch.current) throw new Error("Вход отменён");
    sessionStorage.setItem(TOKEN_KEY, result.access_token);
    const value = await api.me();
    if (version !== epoch.current) throw new Error("Вход отменён");
    setUser(value);
    setError("");
  };
  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        error,
        retry: () => setAttempt((value) => value + 1),
        login,
        register: async (email, password) => {
          await api.register({ email, password });
          await login(email, password);
        },
        logout: clear,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) throw new Error("useAuth must be used inside AuthProvider");
  return value;
}
