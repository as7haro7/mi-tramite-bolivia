import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { MessageSquareWarning, CheckCircle2, XCircle } from 'lucide-react';
import { Modal } from '../components/Modal';

export function ReportesList() {
  const [reportes, setReportes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [selectedReporte, setSelectedReporte] = useState(null);
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    cargarDatos();
  }, []);

  const cargarDatos = async () => {
    try {
      setLoading(true);
      const resp = await api.get('/admin/reportes');
      setReportes(resp.datos || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateStatus = async (estado) => {
    if (!selectedReporte) return;
    setActionLoading(true);
    try {
      await api.put(`/admin/reportes/${selectedReporte.id}`, { estado });
      setSelectedReporte(null);
      cargarDatos();
    } catch (err) {
      alert(err.data?.message || 'Error al actualizar reporte');
    } finally {
      setActionLoading(false);
    }
  };

  if (loading && reportes.length === 0) return <div className="page-loader"><MessageSquareWarning className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;

  return (
    <div className="reportes-list">
      <header className="page-header">
        <h1 className="page-title">Reportes Ciudadanos</h1>
        <p className="page-subtitle">Retroalimentación de los usuarios sobre la información publicada</p>
      </header>

      <div className="card mt-6">
        <div className="table-responsive">
          <table className="table">
            <thead>
              <tr>
                <th>Fecha</th>
                <th>Trámite / Oficina</th>
                <th>Tipo de Error</th>
                <th>Estado</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {reportes.length === 0 ? (
                <tr>
                  <td colSpan="5" className="text-center py-8 text-neutral-500">
                    No hay reportes nuevos. ¡Excelente!
                  </td>
                </tr>
              ) : (
                reportes.map(r => (
                  <tr key={r.id}>
                    <td className="text-sm whitespace-nowrap">
                      {new Date(r.creado_en).toLocaleDateString()}
                    </td>
                    <td>
                      {r.tramite_codigo && <div className="font-bold text-sm">Trámite: {r.tramite_codigo}</div>}
                      {r.oficina_codigo && <div className="font-bold text-sm text-neutral-600">Oficina: {r.oficina_codigo}</div>}
                    </td>
                    <td className="capitalize text-sm font-medium">
                      {r.tipo.replace(/_/g, ' ')}
                    </td>
                    <td>
                      <span className={`badge ${r.estado === 'nuevo' ? 'badge-warning' : r.estado === 'resuelto' ? 'badge-success' : 'badge-neutral'}`}>
                        {r.estado.toUpperCase()}
                      </span>
                    </td>
                    <td>
                      <button className="btn btn-sm btn-outline" onClick={() => setSelectedReporte(r)}>
                        Ver Detalles
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {selectedReporte && (
        <Modal isOpen={true} onClose={() => setSelectedReporte(null)} title="Detalles del Reporte">
          <div className="flex-col gap-4">
            <div className="bg-neutral-50 p-4 rounded-md border border-neutral-200">
              <h4 className="font-bold text-sm text-neutral-500 uppercase tracking-wider mb-2">Mensaje del Ciudadano</h4>
              <p className="text-neutral-900">{selectedReporte.descripcion}</p>
            </div>
            
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-neutral-500 block">Tipo:</span>
                <span className="font-medium capitalize">{selectedReporte.tipo.replace(/_/g, ' ')}</span>
              </div>
              <div>
                <span className="text-neutral-500 block">Fecha de Reporte:</span>
                <span className="font-medium">{new Date(selectedReporte.creado_en).toLocaleString()}</span>
              </div>
            </div>

            {selectedReporte.estado === 'nuevo' || selectedReporte.estado === 'en_revision' ? (
              <div className="flex gap-4 mt-6">
                <button 
                  className="btn flex-1 bg-green-600 hover:bg-green-700 text-white" 
                  disabled={actionLoading}
                  onClick={() => handleUpdateStatus('resuelto')}
                >
                  <CheckCircle2 size={18} /> Marcar Resuelto
                </button>
                <button 
                  className="btn flex-1 bg-neutral-200 hover:bg-neutral-300 text-neutral-800" 
                  disabled={actionLoading}
                  onClick={() => handleUpdateStatus('descartado')}
                >
                  <XCircle size={18} /> Descartar
                </button>
              </div>
            ) : (
              <div className="mt-6 text-center text-sm font-medium text-neutral-500">
                Este reporte ya ha sido {selectedReporte.estado}
              </div>
            )}
          </div>
        </Modal>
      )}
    </div>
  );
}
