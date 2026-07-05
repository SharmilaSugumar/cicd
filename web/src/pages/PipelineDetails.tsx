import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { api } from '../api';
import { Play, CheckCircle2, XCircle, Clock, RefreshCw, AlertTriangle, ArrowLeft, Terminal } from 'lucide-react';

export function PipelineDetails() {
  const { id } = useParams();
  const [runs, setRuns] = useState<any[]>([]);
  const [selectedRun, setSelectedRun] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isRetrying, setIsRetrying] = useState(false);

  const fetchRuns = async () => {
    try {
      const res = await api.get(`/v1/pipelines/${id}/runs`);
      const runsData = res.data || [];
      setRuns(runsData);
      
      if (runsData.length > 0) {
        // If we don't have a selected run, or if we want to refresh the selected run
        const runToSelect = selectedRun ? runsData.find((r: any) => r.id === selectedRun.id) || runsData[0] : runsData[0];
        fetchRunDetails(runToSelect.id);
      } else {
        setIsLoading(false);
      }
    } catch (err) {
      console.error(err);
      setIsLoading(false);
    }
  };

  const fetchRunDetails = async (runId: string) => {
    try {
      const res = await api.get(`/v1/pipelines/${id}/runs/${runId}`);
      setSelectedRun(res.data);
    } catch (err) {
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchRuns();
    // Poll every 3 seconds for live updates
    const interval = setInterval(() => {
      fetchRuns();
    }, 3000);
    return () => clearInterval(interval);
  }, [id]);

  const handleRunPipeline = async () => {
    try {
      await api.post(`/v1/pipelines/${id}/run`, {});
      fetchRuns();
    } catch (err) {
      console.error("Failed to run pipeline:", err);
    }
  };

  const handleRetryJob = async (jobId: string) => {
    setIsRetrying(true);
    try {
      await api.post(`/v1/jobs/${jobId}/retry`, {});
      fetchRunDetails(selectedRun.id);
    } catch (err) {
      console.error("Failed to retry job:", err);
    } finally {
      setIsRetrying(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'COMPLETED':
      case 'SUCCESS':
        return <CheckCircle2 className="h-5 w-5 text-emerald-400" />;
      case 'FAILED':
        return <XCircle className="h-5 w-5 text-rose-400" />;
      case 'RUNNING':
        return <RefreshCw className="h-5 w-5 text-blue-400 animate-spin" />;
      default:
        return <Clock className="h-5 w-5 text-slate-400" />;
    }
  };

  if (isLoading && !selectedRun) {
    return <div className="p-8 text-slate-400">Loading pipeline details...</div>;
  }

  return (
    <div className="h-full flex flex-col fade-in relative space-y-4">
      <div className="flex items-center justify-between pb-4 border-b border-slate-800">
        <div className="flex items-center gap-4">
          <Link to="/pipelines" className="text-slate-400 hover:text-white transition-colors">
            <ArrowLeft className="h-5 w-5" />
          </Link>
          <h1 className="text-2xl font-bold text-slate-100 flex items-center gap-2">
            Pipeline Executions
          </h1>
        </div>
        <div className="flex items-center gap-3">
          <button 
            onClick={async () => {
              if (window.confirm('Are you sure you want to delete this pipeline?')) {
                try {
                  await api.delete(`/v1/pipelines/${id}`);
                  window.location.href = '/pipelines';
                } catch (err) {
                  console.error('Failed to delete pipeline', err);
                }
              }
            }}
            className="flex items-center gap-2 bg-rose-500/10 hover:bg-rose-500/20 text-rose-500 px-4 py-2 rounded-lg transition-colors font-medium border border-rose-500/20"
          >
            <XCircle className="h-4 w-4" /> Delete
          </button>
          <button 
            onClick={handleRunPipeline}
            className="flex items-center gap-2 bg-indigo-500 hover:bg-indigo-600 text-white px-4 py-2 rounded-lg transition-colors font-medium shadow-lg shadow-indigo-500/20"
          >
            <Play className="h-4 w-4 fill-current" /> Trigger Run
          </button>
        </div>
      </div>

      <div className="flex flex-1 gap-6 min-h-0 overflow-hidden">
        {/* Sidebar: List of runs */}
        <div className="w-64 flex flex-col gap-2 overflow-y-auto pr-2 border-r border-slate-800">
          <h2 className="text-xs uppercase font-bold text-slate-500 tracking-wider mb-2 pl-2">Recent Runs</h2>
          {runs.length === 0 ? (
            <div className="text-sm text-slate-500 pl-2">No runs yet.</div>
          ) : (
            runs.map((run) => (
              <button
                key={run.id}
                onClick={() => fetchRunDetails(run.id)}
                className={`flex flex-col text-left p-3 rounded-lg border transition-all ${
                  selectedRun?.id === run.id 
                    ? 'bg-slate-800 border-indigo-500/50 shadow-sm' 
                    : 'bg-slate-900/50 border-transparent hover:bg-slate-800 hover:border-slate-700'
                }`}
              >
                <div className="flex items-center justify-between w-full">
                  <span className="font-mono text-xs text-slate-300 truncate w-24" title={run.id}>{run.id.split('-')[0]}...</span>
                  {getStatusIcon(run.status)}
                </div>
                <div className="text-xs text-slate-500 mt-2">
                  {new Date(run.created_at).toLocaleString()}
                </div>
              </button>
            ))
          )}
        </div>

        {/* Main Content: Logs & Details */}
        <div className="flex-1 flex flex-col min-h-0 bg-slate-900/50 rounded-xl border border-slate-800 overflow-hidden">
          {!selectedRun ? (
            <div className="flex-1 flex items-center justify-center text-slate-500">
              Select a run to view details
            </div>
          ) : (
            <div className="flex flex-col h-full">
              {/* Run Header */}
              <div className="p-4 border-b border-slate-800 bg-slate-800/20 flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-medium text-slate-200 flex items-center gap-2">
                    {getStatusIcon(selectedRun.status)} Run Details
                  </h3>
                  <div className="text-sm text-slate-500 font-mono mt-1">{selectedRun.id}</div>
                </div>
                <div className="text-sm text-slate-400">
                  Status: <span className="text-white font-medium">{selectedRun.status}</span>
                </div>
              </div>

              {/* Jobs & Logs */}
              <div className="flex-1 overflow-y-auto p-4 space-y-6">
                {(!selectedRun.Jobs || selectedRun.Jobs.length === 0) ? (
                  <div className="text-slate-500 text-center py-8">Waiting for jobs to be scheduled...</div>
                ) : (
                  selectedRun.Jobs.map((job: any) => (
                    <div key={job.id} className="space-y-3">
                      <div className="flex items-center justify-between">
                        <h4 className="text-sm font-medium text-slate-300 flex items-center gap-2">
                          <Terminal className="h-4 w-4" /> Job: {job.name} ({job.status})
                        </h4>
                        
                        {/* Retry Button for Failed Jobs */}
                        {job.status === 'FAILED' && (
                          <button
                            onClick={() => handleRetryJob(job.id)}
                            disabled={isRetrying}
                            className="text-xs flex items-center gap-1 bg-slate-800 hover:bg-slate-700 text-slate-200 px-3 py-1.5 rounded transition-colors disabled:opacity-50"
                          >
                            <RefreshCw className={`h-3 w-3 ${isRetrying ? 'animate-spin' : ''}`} />
                            Retry Job
                          </button>
                        )}
                      </div>

                      {/* Dead Letter Queue Warning */}
                      {job.dlq && (
                        <div className="bg-rose-500/10 border border-rose-500/20 rounded-lg p-4 flex gap-3">
                          <AlertTriangle className="h-5 w-5 text-rose-400 flex-shrink-0" />
                          <div>
                            <div className="text-rose-400 font-medium text-sm mb-1">Dead Letter Queue: Execution Failed</div>
                            <div className="text-rose-300/80 text-xs font-mono whitespace-pre-wrap">{job.dlq.error_message}</div>
                          </div>
                        </div>
                      )}

                      {/* Logs Viewer */}
                      <div className="bg-black rounded-lg border border-slate-800 p-4 overflow-x-auto font-mono text-xs text-slate-300 leading-relaxed max-h-96 overflow-y-auto scrollbar-thin scrollbar-thumb-slate-700">
                        {(!job.logs || job.logs.length === 0) ? (
                          <div className="text-slate-600 animate-pulse">Waiting for logs...</div>
                        ) : (
                          job.logs.map((log: any) => (
                            <div key={log.id} className="whitespace-pre-wrap pb-2 border-b border-slate-800/50 mb-2 last:border-0 last:mb-0 last:pb-0">
                              {log.message}
                            </div>
                          ))
                        )}
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
