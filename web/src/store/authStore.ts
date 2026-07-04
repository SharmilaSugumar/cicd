import { create } from 'zustand';

interface AuthState {
  isAuthenticated: boolean;
  user: { name: string; email: string; role: string } | null;
  login: (token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  isAuthenticated: !!localStorage.getItem('auth_token'),
  user: localStorage.getItem('auth_token') ? { name: 'Admin User', email: 'admin@forgeflow.com', role: 'admin' } : null,
  login: (token) => {
    localStorage.setItem('auth_token', token);
    set({ isAuthenticated: true, user: { name: 'Admin User', email: 'admin@forgeflow.com', role: 'admin' } });
  },
  logout: () => {
    localStorage.removeItem('auth_token');
    set({ isAuthenticated: false, user: null });
  },
}));
