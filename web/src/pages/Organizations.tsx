import { useState, useEffect } from 'react';
import { api } from '../api';
import { Building2, Plus, X } from 'lucide-react';

interface Org {
  id: string;
  name: string;
  created_at: string;
}

export function Organizations() {
  const [orgs, setOrgs] = useState<Org[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newOrgName, setNewOrgName] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchOrgs = () => {
    setIsLoading(true);
    api.get('/v1/organizations')
      .then(res => {
        setOrgs(res.data || []);
      })
      .catch(console.error)
      .finally(() => setIsLoading(false));
  };

  useEffect(() => {
    fetchOrgs();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newOrgName.trim()) return;
    
    setIsSubmitting(true);
    try {
      await api.post('/v1/organizations', { name: newOrgName });
      setIsModalOpen(false);
      setNewOrgName('');
      fetchOrgs();
    } catch (error) {
      console.error(error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6 fade-in h-full flex flex-col relative">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-slate-100 flex items-center gap-2">
            <Building2 className="h-6 w-6 text-indigo-400" /> Organizations
          </h1>
          <p className="mt-1 text-sm text-slate-400">Manage your tenant organizations and namespaces.</p>
        </div>
        <button 
          onClick={() => setIsModalOpen(true)}
          className="flex items-center gap-2 bg-indigo-500 hover:bg-indigo-600 text-white px-4 py-2 rounded-lg transition-colors font-medium shadow-lg shadow-indigo-500/20"
        >
          <Plus className="h-4 w-4" /> New Organization
        </button>
      </div>

      <div className="rounded-xl border border-slate-800 bg-slate-900 overflow-hidden flex-1 shadow-sm">
        {isLoading ? (
          <div className="p-8 text-center text-slate-400">Loading organizations...</div>
        ) : orgs.length === 0 ? (
          <div className="p-16 text-center text-slate-400">No organizations found.</div>
        ) : (
          <table className="w-full text-left text-sm text-slate-300">
            <thead className="bg-slate-800/50 text-xs uppercase text-slate-400">
              <tr>
                <th className="px-6 py-4 font-medium">Name</th>
                <th className="px-6 py-4 font-medium">ID</th>
                <th className="px-6 py-4 font-medium">Created At</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800">
              {orgs.map(org => (
                <tr key={org.id} className="hover:bg-slate-800/50 transition-colors">
                  <td className="px-6 py-4 font-medium text-slate-200">{org.name}</td>
                  <td className="px-6 py-4 text-slate-500 font-mono text-xs">{org.id}</td>
                  <td className="px-6 py-4 text-slate-500">{new Date(org.created_at).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm transition-opacity fade-in">
          <div className="bg-slate-900 border border-slate-800 rounded-xl shadow-2xl w-full max-w-md overflow-hidden transform transition-all">
            <div className="flex justify-between items-center p-6 border-b border-slate-800/60">
              <h2 className="text-xl font-semibold text-slate-100">Create Organization</h2>
              <button onClick={() => setIsModalOpen(false)} className="text-slate-400 hover:text-white transition-colors">
                <X className="h-5 w-5" />
              </button>
            </div>
            <form onSubmit={handleCreate} className="p-6 space-y-6">
              <div>
                <label className="block text-sm font-medium text-slate-300 mb-2">Organization Name</label>
                <input
                  type="text"
                  value={newOrgName}
                  onChange={(e) => setNewOrgName(e.target.value)}
                  className="w-full bg-slate-950/50 border border-slate-800 rounded-lg px-4 py-2.5 text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all placeholder:text-slate-600"
                  placeholder="e.g. Acme Corp"
                  autoFocus
                  required
                />
              </div>
              <div className="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 text-sm font-medium text-slate-300 hover:text-white hover:bg-slate-800 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting || !newOrgName.trim()}
                  className="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white text-sm font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-indigo-500/20"
                >
                  {isSubmitting ? 'Creating...' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
