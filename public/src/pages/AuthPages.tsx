import { Eye, EyeOff, Leaf } from "lucide-react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useState } from "react";
import { useAuth } from "../auth";
import { ApiError } from "../api";
import { Button, Field, Notice } from "../shared/ui";

function AuthFrame({
  title,
  subtitle,
  children,
  footer,
}: {
  title: string;
  subtitle: string;
  children: React.ReactNode;
  footer: React.ReactNode;
}) {
  return (
    <div className="auth-page">
      <div className="auth-decoration">
        <span className="circle c1" />
        <span className="circle c2" />
        <div className="auth-quote">
          <Leaf size={38} />
          <p>
            Сохраняйте истории,
            <br />
            которые связывают поколения.
          </p>
        </div>
      </div>
      <div className="auth-card">
        <div className="brand auth-brand">
          <Leaf size={27} />
          <div>
            <strong>Roots</strong>
            <small>Your family story</small>
          </div>
        </div>
        <h1>{title}</h1>
        <p className="auth-subtitle">{subtitle}</p>
        {children}
        <div className="auth-footer">{footer}</div>
      </div>
    </div>
  );
}
function AuthForm({ register = false }: { register?: boolean }) {
  const { login, register: signUp } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [show, setShow] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const returnTo = new URLSearchParams(location.search).get("returnTo");
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (loading) return;
    setError("");
    if (register && password !== confirm) {
      setError("Пароли не совпадают");
      return;
    }
    const passwordBytes = new TextEncoder().encode(password).length;
    if (passwordBytes < 8 || (register && passwordBytes > 72)) {
      setError("Пароль должен содержать от 8 до 72 байт");
      return;
    }
    setLoading(true);
    try {
      if (register) await signUp(email.trim(), password);
      else await login(email.trim(), password);
      navigate(
        returnTo?.startsWith("/") &&
          !returnTo.startsWith("//") &&
          !returnTo.includes("\\")
          ? returnTo
          : "/trees",
        {
          replace: true,
        },
      );
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Не удалось выполнить запрос",
      );
    } finally {
      setLoading(false);
    }
  };
  return (
    <form className="auth-form" onSubmit={submit}>
      {error && <Notice>{error}</Notice>}
      <Field label="Email">
        <input
          type="email"
          autoComplete="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="you@example.com"
          required
        />
      </Field>
      <Field label="Пароль" hint={register ? "От 8 до 72 байт" : undefined}>
        <div className="password-field">
          <input
            type={show ? "text" : "password"}
            autoComplete={register ? "new-password" : "current-password"}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Введите пароль"
            minLength={8}
            required
          />
          <button
            type="button"
            onClick={() => setShow(!show)}
            aria-label={show ? "Скрыть пароль" : "Показать пароль"}
          >
            {show ? <EyeOff size={17} /> : <Eye size={17} />}
          </button>
        </div>
      </Field>
      {register && (
        <Field label="Подтвердите пароль">
          <input
            type={show ? "text" : "password"}
            autoComplete="new-password"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            placeholder="Повторите пароль"
            minLength={8}
            required
          />
        </Field>
      )}
      <Button type="submit" loading={loading}>
        {register ? "Создать аккаунт" : "Войти"}
      </Button>
    </form>
  );
}
export function LoginPage() {
  return (
    <AuthFrame
      title="С возвращением"
      subtitle="Войдите, чтобы продолжить историю своей семьи."
      footer={
        <>
          Нет аккаунта? <Link to="/register">Создать аккаунт</Link>
        </>
      }
    >
      <AuthForm />
    </AuthFrame>
  );
}
export function RegisterPage() {
  return (
    <AuthFrame
      title="Начните свою историю"
      subtitle="Создайте аккаунт и сохраните важные связи."
      footer={
        <>
          Уже есть аккаунт? <Link to="/login">Войти</Link>
        </>
      }
    >
      <AuthForm register />
    </AuthFrame>
  );
}
