import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { 
  FileEdit, 
  AlertTriangle, 
  DatabaseBackup, 
  MessageSquareWarning,
  Activity,
  Network
} from 'lucide-react';

export function Dashboard() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    cargarDashboard();
  }, []);

  const cargarDashboard = async () => {
    try {
      setLoading(true);
      const resp = await api.get('/admin/dashboard');
      setData(resp);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="page-loader"><Activity className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;
  if (!data) return null;

  const m = data.metricas;

  return (
    <div className="dashboard">
      <header className="page-header">
        <div>
          <h1 className="page-title">Dashboard Operativo</h1>
          <p className="page-subtitle">Métricas de la plataforma v2</p>
        </div>
      </header>

      <div className="metrics-grid">
        <div className="metric-card warning">
          <div className="metric-icon"><FileEdit size={24} /></div>
          <div className="metric-content">
            <span className="metric-value">{m.versiones_en_revision}</span>
            <span className="metric-label">Versiones en Revisión</span>
          </div>
        </div>

        <div className="metric-card danger">
          <div className="metric-icon"><AlertTriangle size={24} /></div>
          <div className="metric-content">
            <span className="metric-value">{m.tramites_revision_vencida}</span>
            <span className="metric-label">Revisiones Vencidas</span>
          </div>
        </div>

        <div className="metric-card danger">
          <div className="metric-icon"><DatabaseBackup size={24} /></div>
          <div className="metric-content">
            <span className="metric-value">{m.embeddings_fallidos}</span>
            <span className="metric-label">Embeddings Fallidos</span>
          </div>
        </div>

        <div className="metric-card warning">
          <div className="metric-icon"><MessageSquareWarning size={24} /></div>
          <div className="metric-content">
            <span className="metric-value">{m.reportes_ciudadanos_nuevos}</span>
            <span className="metric-label">Reportes Ciudadanos</span>
          </div>
        </div>

        <div className="metric-card primary">
          <div className="metric-icon"><Network size={24} /></div>
          <div className="metric-content">
            <span className="metric-value">{m.candidatos_pendientes}</span>
            <span className="metric-label">Candidatos de Ingesta</span>
          </div>
        </div>
      </div>

      <div className="dashboard-section mt-8">
        <h2 className="section-title">Publicaciones por Institución</h2>
        <div className="card">
          {data.por_institucion?.length === 0 ? (
            <p className="empty-state">No hay publicaciones activas.</p>
          ) : (
            <div className="list-group">
              {data.por_institucion.map((inst, idx) => (
                <div key={idx} className="list-item flex justify-between">
                  <span>{inst.institucion}</span>
                  <span className="badge badge-primary">{inst.publicaciones} publicados</span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
