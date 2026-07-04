import { useState, useEffect } from 'react';
import { Activity, GitMerge, FolderKanban, Users } from 'lucide-react';
import { api } from '../api';

export function Dashboard() {
  const [pipelineCount, setPipelineCount] = useState(0);
  const [projectCount, setProjectCount] = useState(0);
  const [orgCount, setOrgCount] = useState(0);

  useEffect(() => {
    Promise.all([
      api.get('/v1/pipelines').catch(() => ({ data: [] })),
      api.get('/v1/projects').catch(() => ({ data: [] })),
      api.get('/v1/organizations').catch(() => ({ data: [] }))
    ]).then(([pipelines, projects, orgs]) => {
      setPipelineCount(pipelines.data?.length || 0);
      setProjectCount(projects.data?.length || 0);
      setOrgCount(orgs.data?.length || 0);
    }).catch(console.error);
  }, []);

  const stats = [
    { name: 'Total Pipelines', value: pipelineCount.toString(), change: '+12.5%', icon: GitMerge, trend: 'up' },
    { name: 'Total Projects', value: projectCount.toString(), change: '+4.1%', icon: FolderKanban, trend: 'up' },
    { name: 'Total Organizations', value: orgCount.toString(), change: '0%', icon: Users, trend: 'neutral' },
    { name: 'Success Rate', value: '98.4%', change: '+0.4%', icon: Activity, trend: 'up' },
  ];

  return (
    <div className="space-y-6 fade-in">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-slate-100">Dashboard</h1>
          <p className="mt-1 text-sm text-slate-400">Overview of your ForgeFlow instance</p>
        </div>
        <button className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500 transition-colors shadow-lg shadow-indigo-500/20">
          Create Pipeline
        </button>
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <div
            key={stat.name}
            className="rounded-xl border border-slate-800 bg-slate-900/50 p-6 shadow-sm backdrop-blur-sm transition-all hover:bg-slate-800/50 hover:border-slate-700"
          >
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium text-slate-400">{stat.name}</p>
              <div className="rounded-md bg-slate-800 p-2 ring-1 ring-white/10">
                <stat.icon className="h-5 w-5 text-indigo-400" />
              </div>
            </div>
            <div className="mt-4 flex items-baseline gap-2">
              <p className="text-3xl font-semibold text-slate-100">{stat.value}</p>
              <p className={`text-sm font-medium ${
                stat.trend === 'up' ? 'text-emerald-400' : stat.trend === 'down' ? 'text-rose-400' : 'text-slate-500'
              }`}>
                {stat.change}
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
