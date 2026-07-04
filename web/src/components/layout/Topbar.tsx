import { Bell, Search, User, Moon, Menu } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';

export function Topbar() {
  const { user, logout } = useAuthStore();

  return (
    <header className="flex h-16 items-center justify-between border-b border-slate-800 bg-slate-900/80 px-6 backdrop-blur-md sticky top-0 z-50">
      <div className="flex items-center gap-4">
        <button className="lg:hidden text-slate-400 hover:text-slate-100 transition-colors">
          <Menu className="h-6 w-6" />
        </button>
        <div className="relative hidden md:block">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500" />
          <input
            type="text"
            placeholder="Search anywhere..."
            className="h-10 w-80 rounded-full border border-slate-700 bg-slate-800/50 pl-10 pr-4 text-sm text-slate-200 placeholder:text-slate-500 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 transition-all duration-200"
          />
        </div>
      </div>

      <div className="flex items-center gap-5">
        <div className="flex items-center gap-3 border-r border-slate-700 pr-5">
          <button className="relative text-slate-400 hover:text-slate-100 transition-colors">
            <Bell className="h-5 w-5" />
            <span className="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-indigo-500 text-[10px] font-bold text-white shadow-sm ring-2 ring-slate-900">
              3
            </span>
          </button>
          <button className="text-slate-400 hover:text-slate-100 transition-colors">
            <Moon className="h-5 w-5" />
          </button>
        </div>

        <div className="flex items-center gap-3">
          <div className="flex flex-col items-end hidden sm:flex">
            <span className="text-sm font-medium text-slate-200">{user?.name || 'User'}</span>
            <span className="text-xs text-slate-500 capitalize">{user?.role || 'Guest'}</span>
          </div>
          <button 
            className="flex h-9 w-9 items-center justify-center rounded-full bg-slate-800 ring-2 ring-transparent hover:ring-indigo-500 transition-all duration-200 focus:outline-none"
            onClick={() => {
              if (window.confirm('Are you sure you want to log out?')) {
                logout();
              }
            }}
            title="Click to logout"
          >
            <User className="h-5 w-5 text-slate-400" />
          </button>
        </div>
      </div>
    </header>
  );
}
