import { Navigate, Route, Routes, useLocation } from "react-router-dom";
import { useAuth } from "../auth";
import { LoginPage, RegisterPage } from "../pages/AuthPages";
import { Shell } from "../shared/Shell";
import { lazy, Suspense } from "react";
const TreePage = lazy(() =>
  import("../pages/TreePage").then((module) => ({ default: module.TreePage })),
);
const FamiliesPage = lazy(() =>
  import("../pages/FamiliesPage").then((module) => ({
    default: module.FamiliesPage,
  })),
);
const SettingsPage = lazy(() =>
  import("../pages/SettingsPage").then((module) => ({
    default: module.SettingsPage,
  })),
);

function Protected({ children }: { children: React.ReactNode }) {
  const { user, loading, error, retry, logout } = useAuth();
  const location = useLocation();
  if (loading)
    return (
      <div className="page-loader">
        <span className="spinner" />
        Загрузка Roots
      </div>
    );
  if (error)
    return (
      <div className="center-panel" role="alert">
        <p>{error}</p>
        <button className="button" onClick={retry}>
          Повторить
        </button>
        <button className="button secondary" onClick={logout}>
          Выйти
        </button>
      </div>
    );
  return user ? (
    <>{children}</>
  ) : (
    <Navigate
      to={`/login?returnTo=${encodeURIComponent(location.pathname + location.search)}`}
      replace
    />
  );
}
export function App() {
  return (
    <Suspense fallback={<div className="page-loader">Загрузка раздела…</div>}>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route
          element={
            <Protected>
              <Shell />
            </Protected>
          }
        >
          <Route path="/" element={<Navigate to="/trees" replace />} />
          <Route path="/trees" element={<TreePage />} />
          <Route path="/trees/:treeId" element={<TreePage />} />
          <Route path="/trees/:treeId/families" element={<FamiliesPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Suspense>
  );
}
