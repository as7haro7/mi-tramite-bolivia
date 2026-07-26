import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { History, Eye } from 'lucide-react';
import { Modal } from '../components/Modal';

export function AuditoriaList() {
  const [eventos, setEventos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [selectedEvento, setSelectedEvento] = useState(null);

  useEffect(() => {
    cargarDatos();
  }, []);

  const cargarDatos = async () => {
    try {
      setLoading(true);
      const resp = await api.get('/admin/auditoria');
      setEventos(resp.datos || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading && eventos.length === 0) return <div className="page-loader"><History className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;

  return (
    <div className="auditoria-list">
      <header className="page-header">
        <h1 className="page-title">Registro de Auditoría</h1>
        <p className="page-subtitle">Trazabilidad de acciones administrativas en el sistema</p>
      </header>

      <div className="card mt-6">
        <div className="table-responsive">
          <table className="table">
            <thead>
              <tr>
                <th>Fecha y Hora</th>
                <th>Usuario</th>
                <th>Acción</th>
                <th>Entidad Afectada</th>
                <th>IP</th>
                <th>Detalles</th>
              </tr>
            </thead>
            <tbody>
              {eventos.length === 0 ? (
                <tr>
                  <td colSpan="6" className="text-center py-8 text-neutral-500">
                    No hay eventos de auditoría registrados.
                  </td>
                </tr>
              ) : (
                eventos.map(e => (
                  <tr key={e.id}>
                    <td className="text-sm whitespace-nowrap">
                      {new Date(e.ocurrido_en).toLocaleString()}
                    </td>
                    <td className="font-medium text-sm">{e.actor_correo || 'Sistema'}</td>
                    <td className="font-bold text-sm text-primary-600">{e.accion}</td>
                    <td>
                      <span className="text-sm font-mono text-neutral-600">{e.entidad_tipo}</span>
                      <div className="text-xs text-neutral-400">ID: {e.entidad_id}</div>
                    </td>
                    <td className="text-xs font-mono">{e.ip_hash?.substring(0, 8) || '-'}</td>
                    <td>
                      <button className="btn btn-sm btn-outline" onClick={() => setSelectedEvento(e)}>
                        <Eye size={14} />
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {selectedEvento && (
        <Modal isOpen={true} onClose={() => setSelectedEvento(null)} title="Detalle de Auditoría">
          <div className="flex-col gap-4">
            <div className="grid grid-cols-2 gap-4 text-sm bg-neutral-50 p-4 rounded-md">
              <div><span className="text-neutral-500 block">Acción:</span> <span className="font-bold">{selectedEvento.accion}</span></div>
              <div><span className="text-neutral-500 block">Usuario:</span> {selectedEvento.actor_correo}</div>
              <div><span className="text-neutral-500 block">Fecha:</span> {new Date(selectedEvento.ocurrido_en).toLocaleString()}</div>
              <div><span className="text-neutral-500 block">Entidad:</span> {selectedEvento.entidad_tipo} ({selectedEvento.entidad_id})</div>
              <div className="col-span-2"><span className="text-neutral-500 block">User Agent:</span> <span className="font-mono text-xs">{selectedEvento.user_agent}</span></div>
            </div>

            <div className="grid grid-cols-2 gap-4 mt-2">
              <div>
                <h4 className="font-bold text-xs text-neutral-500 uppercase mb-2">Estado Anterior</h4>
                <pre className="bg-neutral-900 text-neutral-100 p-2 rounded text-xs overflow-auto max-h-64">
                  {selectedEvento.antes ? JSON.stringify(selectedEvento.antes, null, 2) : 'Ninguno'}
                </pre>
              </div>
              <div>
                <h4 className="font-bold text-xs text-neutral-500 uppercase mb-2">Nuevo Estado</h4>
                <pre className="bg-neutral-900 text-neutral-100 p-2 rounded text-xs overflow-auto max-h-64">
                  {selectedEvento.despues ? JSON.stringify(selectedEvento.despues, null, 2) : 'Ninguno'}
                </pre>
              </div>
            </div>
          </div>
        </Modal>
      )}
    </div>
  );
}
