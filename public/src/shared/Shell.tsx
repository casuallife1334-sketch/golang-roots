import { NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";
import { GitBranch, LogOut, Settings, UsersRound, Leaf } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { api } from "../api";
import { useAuth } from "../auth";

export function Shell() {
  const { user, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const trees = useQuery({ queryKey: ["trees", user?.id], queryFn: api.trees });
  const routeTreeId = location.pathname.match(/^\/trees\/([^/]+)/)?.[1];
  const storedTreeId = user
    ? sessionStorage.getItem(`roots:last-tree:${user.id}`)
    : null;
  const activeTreeId =
    routeTreeId ||
    (storedTreeId && trees.data?.some((tree) => tree.id === storedTreeId)
      ? storedTreeId
      : trees.data?.[0]?.id);
  const treeHref = activeTreeId ? `/trees/${activeTreeId}` : "/trees";
  const familiesHref = activeTreeId
    ? `/trees/${activeTreeId}/families`
    : "/trees";
  const initials = user?.email.slice(0, 2).toUpperCase() || "R";
  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <Leaf size={28} strokeWidth={1.7} />
          <div>
            <strong>Roots</strong>
            <small>Your family story</small>
          </div>
        </div>
        <nav className="main-nav">
          <NavLink
            to={treeHref}
            end
            className={() =>
              `nav-item ${location.pathname === treeHref ? "active" : ""}`
            }
          >
            <GitBranch size={19} />
            Древо
          </NavLink>
          <NavLink
            to={familiesHref}
            end
            className={() =>
              `nav-item ${location.pathname === familiesHref ? "active" : ""}`
            }
          >
            <UsersRound size={19} />
            Семьи
          </NavLink>
          <NavLink
            to="/settings"
            className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}
          >
            <Settings size={19} />
            Настройки
          </NavLink>
        </nav>
        <div className="sidebar-bottom">
          <button className="user-chip" onClick={() => navigate("/settings")}>
            <span className="avatar small">{initials}</span>
            <span>{user?.email}</span>
          </button>
          <button
            className="icon-button"
            aria-label="Выйти"
            title="Выйти"
            onClick={logout}
          >
            <LogOut size={17} />
          </button>
        </div>
      </aside>
      <main className="main-area">
        <Outlet />
      </main>
    </div>
  );
}
