import { Navigate, Route, Routes, useLocation } from 'react-router-dom'
import { useAuth } from '../auth'
import { LoginPage, RegisterPage } from '../pages/AuthPages'
import { Shell } from '../shared/Shell'
import { TreePage } from '../pages/TreePage'
import { FamiliesPage } from '../pages/FamiliesPage'
import { SettingsPage } from '../pages/SettingsPage'

function Protected({ children }: { children: React.ReactNode }) { const { user, loading } = useAuth(); const location = useLocation(); if (loading) return <div className="page-loader"><span className="spinner" />Загрузка Roots</div>; return user ? <>{children}</> : <Navigate to={`/login?returnTo=${encodeURIComponent(location.pathname + location.search)}`} replace /> }
export function App() { return <Routes><Route path="/login" element={<LoginPage />} /><Route path="/register" element={<RegisterPage />} /><Route element={<Protected><Shell /></Protected>}><Route path="/" element={<Navigate to="/trees" replace />} /><Route path="/trees" element={<TreePage />} /><Route path="/trees/:treeId" element={<TreePage />} /><Route path="/trees/:treeId/families" element={<FamiliesPage />} /><Route path="/settings" element={<SettingsPage />} /></Route><Route path="*" element={<Navigate to="/" replace />} /></Routes> }
