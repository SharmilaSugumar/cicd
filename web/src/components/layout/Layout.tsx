import { Outlet, Navigate } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { Topbar } from './Topbar';
import { useAuthStore } from '../../store/authStore';

export function Layout() {
  const { isAuthenticated } = useAuthStore();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="flex h-screen w-full bg-slate-950 text-slate-200 font-sans selection:bg-indigo-500/30">
      <Sidebar />
      <div className="flex flex-1 flex-col overflow-hidden relative">
        <Topbar />
        <main className="flex-1 overflow-x-hidden overflow-y-auto bg-slate-950/50 p-6">
          <div className="mx-auto max-w-7xl h-full">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}
