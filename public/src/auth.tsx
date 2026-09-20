import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { useNavigate } from "react-router-dom";
import { api, setUnauthorizedHandler } from "./api";
import type { User } from "./types";

interface AuthContextValue {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
}
const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const clear = () => {
    sessionStorage.removeItem("roots:access-token");
    setUser(null);
    sessionStorage.setItem(
      "roots:session-message",
      "Сессия истекла. Войдите снова",
    );
    navigate("/login", { replace: true });
  };
  useEffect(() => {
    setUnauthorizedHandler(clear);
    const existing = sessionStorage.getItem("roots:access-token");
    if (!existing) {
      setLoading(false);
      return;
    }
    api
      .me()
      .then(setUser)
      .catch(clear)
      .finally(() => setLoading(false));
  }, []);
  const login = async (email: string, password: string) => {
    const result = await api.login({ email, password });
    sessionStorage.setItem("roots:access-token", result.access_token);
    setUser(await api.me());
  };
  const register = async (email: string, password: string) => {
    await api.register({ email, password });
    await login(email, password);
  };
  const logout = () => {
    sessionStorage.removeItem("roots:access-token");
    setUser(null);
    navigate("/login", { replace: true });
  };
  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) throw new Error("useAuth must be used inside AuthProvider");
  return value;
}
