import { useState, useEffect } from 'react';
import { Activity, GitMerge, FolderKanban, Users } from 'lucide-react';
import { api } from '../api';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar, Cell } from 'recharts';



export function Dashboard() {
  const [metrics, setMetrics] = useState<any>({
    pipeline_count: 0,
    project_count: 0,
    org_count: 0,
    success_rate: 0,
    activity_data: [],
    language_data: []
  });

  useEffect(() => {
    api.get('/v1/metrics')
      .then(res => {
        setMetrics(res.data);
      })
      .catch(console.error);
  }, []);

  const stats = [
    { name: 'Total Pipelines', value: metrics.pipeline_count.toString(), change: '+12.5%', icon: GitMerge, trend: 'up' },
    { name: 'Total Projects', value: metrics.project_count.toString(), change: '+4.1%', icon: FolderKanban, trend: 'up' },
    { name: 'Total Organizations', value: metrics.org_count.toString(), change: '0%', icon: Users, trend: 'neutral' },
    { name: 'Success Rate', value: `${metrics.success_rate.toFixed(1)}%`, change: '+0.4%', icon: Activity, trend: 'up' },
  ];

  return (
    <div className="space-y-6 fade-in h-full flex flex-col">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-slate-100">Dashboard</h1>
          <p className="mt-1 text-sm text-slate-400">Overview of your ForgeFlow instance</p>
        </div>
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

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 flex-1 min-h-0">
        <div className="lg:col-span-2 rounded-xl border border-slate-800 bg-slate-900/50 p-6 shadow-sm backdrop-blur-sm flex flex-col">
          <h2 className="text-base font-semibold text-slate-100 mb-6">Pipeline Activity (7 Days)</h2>
          <div className="flex-1 min-h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={metrics.activity_data} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                <defs>
                  <linearGradient id="colorPipelines" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#818cf8" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#818cf8" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorSuccess" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#34d399" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#34d399" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} />
                <XAxis dataKey="name" stroke="#64748b" tick={{fill: '#64748b', fontSize: 12}} axisLine={false} tickLine={false} dy={10} />
                <YAxis stroke="#64748b" tick={{fill: '#64748b', fontSize: 12}} axisLine={false} tickLine={false} />
                <Tooltip 
                  contentStyle={{ backgroundColor: '#0f172a', borderColor: '#1e293b', borderRadius: '0.5rem', color: '#f1f5f9' }}
                  itemStyle={{ color: '#f1f5f9' }}
                />
                <Area type="monotone" dataKey="pipelines" name="Total Runs" stroke="#818cf8" strokeWidth={2} fillOpacity={1} fill="url(#colorPipelines)" />
                <Area type="monotone" dataKey="success" name="Successful" stroke="#34d399" strokeWidth={2} fillOpacity={1} fill="url(#colorSuccess)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-6 shadow-sm backdrop-blur-sm flex flex-col">
          <h2 className="text-base font-semibold text-slate-100 mb-6">Language Distribution</h2>
          <div className="flex-1 min-h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={metrics.language_data} layout="vertical" margin={{ top: 0, right: 30, left: 10, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" horizontal={true} vertical={false} />
                <XAxis type="number" hide />
                <YAxis dataKey="name" type="category" stroke="#64748b" tick={{fill: '#64748b', fontSize: 13}} axisLine={false} tickLine={false} />
                <Tooltip 
                  cursor={{fill: '#1e293b', opacity: 0.4}}
                  contentStyle={{ backgroundColor: '#0f172a', borderColor: '#1e293b', borderRadius: '0.5rem', color: '#f1f5f9' }}
                />
                <Bar dataKey="count" radius={[0, 4, 4, 0]} barSize={24}>
                  {metrics.language_data.map((entry: any, index: number) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
    </div>
  );
}
