import { Link, useLocation } from 'react-router-dom';
import { cn } from '../../lib/utils';
import {
  LayoutDashboard,
  Building2,
  FolderKanban,
  GitMerge
} from 'lucide-react';

const navItems = [
  { name: 'Dashboard', path: '/', icon: LayoutDashboard },
  { name: 'Organizations', path: '/organizations', icon: Building2 },
  { name: 'Projects', path: '/projects', icon: FolderKanban },
  { name: 'Pipelines', path: '/pipelines', icon: GitMerge },

];

export function Sidebar() {
  const location = useLocation();

  return (
    <div className="flex h-full w-64 flex-col bg-slate-900 border-r border-slate-800 text-slate-300 transition-all duration-300">
      <div className="flex h-16 items-center px-6 border-b border-slate-800 bg-slate-950/50">
        <div className="flex items-center gap-3">
          <div className="h-8 w-8 rounded bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center font-bold text-white shadow-lg shadow-indigo-500/20">
            FF
          </div>
          <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-white to-slate-400">ForgeFlow</span>
        </div>
      </div>
      
      <div className="flex-1 overflow-y-auto py-4 scrollbar-thin scrollbar-thumb-slate-700">
        <nav className="space-y-1 px-3">
          {navItems.map((item) => {
            const isActive = location.pathname === item.path || (item.path !== '/' && location.pathname.startsWith(item.path));
            return (
              <Link
                key={item.name}
                to={item.path}
                className={cn(
                  "group flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200",
                  isActive
                    ? "bg-indigo-500/10 text-indigo-400"
                    : "hover:bg-slate-800/50 hover:text-slate-100"
                )}
              >
                <item.icon className={cn(
                  "h-5 w-5 transition-colors duration-200",
                  isActive ? "text-indigo-400" : "text-slate-500 group-hover:text-slate-300"
                )} />
                {item.name}
              </Link>
            );
          })}
        </nav>
      </div>
      
      <div className="p-4 border-t border-slate-800">
        <div className="rounded-xl bg-slate-800/50 p-4 ring-1 ring-white/10">
          <p className="text-xs font-semibold text-slate-400">System Status</p>
          <div className="mt-2 flex items-center gap-2">
            <div className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse"></div>
            <p className="text-sm font-medium text-slate-200">All systems operational</p>
          </div>
        </div>
      </div>
    </div>
  );
}
