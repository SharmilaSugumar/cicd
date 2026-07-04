

export function GenericPage({ title, description }: { title: string, description: string }) {
  return (
    <div className="space-y-6 fade-in h-full flex flex-col">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-slate-100">{title}</h1>
          <p className="mt-1 text-sm text-slate-400">{description}</p>
        </div>
        <button className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500 transition-colors shadow-lg shadow-indigo-500/20">
          Create New
        </button>
      </div>
      
      <div className="flex-1 rounded-xl border border-slate-800 bg-slate-900/50 p-8 shadow-sm flex flex-col items-center justify-center text-center">
        <div className="w-16 h-16 rounded-full bg-slate-800 flex items-center justify-center mb-4 ring-1 ring-white/10">
          <svg className="w-8 h-8 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-slate-200 mb-2">No {title.toLowerCase()} found</h3>
        <p className="text-slate-400 text-sm max-w-sm mb-6">
          Get started by creating your first {title.toLowerCase().replace(/s$/, '')}. It only takes a few minutes to set up.
        </p>
        <button className="rounded-lg border border-slate-700 bg-slate-800 px-4 py-2 text-sm font-medium text-slate-200 hover:bg-slate-700 hover:text-white transition-colors">
          Learn More
        </button>
      </div>
    </div>
  );
}
