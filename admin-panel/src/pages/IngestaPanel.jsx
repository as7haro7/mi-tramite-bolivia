import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { Database, Link2, Check, X, RefreshCw } from 'lucide-react';
import { Modal } from '../components/Modal';

export function IngestaPanel() {
  const [fuentes, setFuentes] = useState([]);
  const [candidatos, setCandidatos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [activeTab, setActiveTab] = useState('fuentes'); // 'fuentes' | 'candidatos'
  
  const [selectedCandidato, setSelectedCandidato] = useState(null);
  const [actionLoading, setActionLoading] = useState(false);

  const [isFuenteModalOpen, setIsFuenteModalOpen] = useState(false);
  const [isEditingFuente, setIsEditingFuente] = useState(false);
  const [fuenteForm, setFuenteForm] = useState({ 
    id: null,
    nombre: '', 
    tipo: 'web', 
    url: '', 
    institucion_id: '',
    frecuencia_cron: '',
    estado: 'activa'
  });
  const [fuenteLoading, setFuenteLoading] = useState(false);

  useEffect(() => {
    cargarDatos();
  }, []);

  const cargarDatos = async () => {
    try {
      setLoading(true);
      const [fResp, cResp] = await Promise.all([
        api.get('/admin/fuentes'),
        api.get('/admin/candidatos')
      ]);
      setFuentes(fResp.datos || []);
      setCandidatos(cResp.datos || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateCandidato = async (estado) => {
    if (!selectedCandidato) return;
    setActionLoading(true);
    try {
      await api.put(`/admin/candidatos/${selectedCandidato.id}`, { estado });
      setSelectedCandidato(null);
      cargarDatos();
    } catch (err) {
      alert(err.data?.message || 'Error al actualizar candidato');
    } finally {
      setActionLoading(false);
    }
  };

  const handleCreateFuente = async (e) => {
    e.preventDefault();
    setFuenteLoading(true);
    try {
      const payload = {
        ...fuenteForm,
        institucion_id: fuenteForm.institucion_id ? parseInt(fuenteForm.institucion_id, 10) : undefined,
        frecuencia_cron: fuenteForm.frecuencia_cron || undefined,
      };

      if (isEditingFuente) {
        await api.put(`/admin/fuentes/${fuenteForm.id}`, payload);
      } else {
        await api.post('/admin/fuentes', payload);
      }

      setIsFuenteModalOpen(false);
      setFuenteForm({ id: null, nombre: '', tipo: 'web', url: '', institucion_id: '', frecuencia_cron: '', estado: 'activa' });
      cargarDatos();
    } catch (err) {
      alert(err.data?.message || 'Error al guardar fuente');
    } finally {
      setFuenteLoading(false);
    }
  };

  const handleDeleteFuente = async (id) => {
    if (!confirm('¿Seguro que deseas eliminar esta fuente?')) return;
    try {
      await api.delete(`/admin/fuentes/${id}`);
      cargarDatos();
    } catch (err) {
      alert(err.data?.message || 'Error al eliminar');
    }
  };

  const openNewFuente = () => {
    setIsEditingFuente(false);
    setFuenteForm({ id: null, nombre: '', tipo: 'web', url: '', institucion_id: '', frecuencia_cron: '', estado: 'activa' });
    setIsFuenteModalOpen(true);
  };

  const openEditFuente = (f) => {
    setIsEditingFuente(true);
    setFuenteForm({
      id: f.id,
      nombre: f.nombre,
      tipo: f.tipo,
      url: f.url,
      institucion_id: f.institucion_id || '',
      frecuencia_cron: f.frecuencia_cron || '',
      estado: f.estado || 'activa'
    });
    setIsFuenteModalOpen(true);
  };

  if (loading && fuentes.length === 0) return <div className="page-loader"><Database className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;

  return (
    <div className="ingesta-panel">
      <header className="page-header flex justify-between items-center">
        <div>
          <h1 className="page-title">Ingesta y Scraping (IA)</h1>
          <p className="page-subtitle">Gestión de fuentes oficiales y revisión de contenido candidato extraído</p>
        </div>
        {activeTab === 'fuentes' && (
          <button className="btn btn-primary" onClick={openNewFuente}>
            Nueva Fuente
          </button>
        )}
      </header>

      <div className="flex gap-4 mt-6 border-b border-neutral-200">
        <button 
          className={`pb-2 px-4 font-bold border-b-2 ${activeTab === 'fuentes' ? 'border-primary-500 text-primary-600' : 'border-transparent text-neutral-500 hover:text-neutral-700'}`}
          onClick={() => setActiveTab('fuentes')}
        >
          Fuentes de Origen ({fuentes.length})
        </button>
        <button 
          className={`pb-2 px-4 font-bold border-b-2 ${activeTab === 'candidatos' ? 'border-primary-500 text-primary-600' : 'border-transparent text-neutral-500 hover:text-neutral-700'}`}
          onClick={() => setActiveTab('candidatos')}
        >
          Candidatos Pendientes ({candidatos.filter(c => c.estado === 'pendiente').length})
        </button>
      </div>

      <div className="card mt-6">
        {activeTab === 'fuentes' && (
          <div className="table-responsive">
            <table className="table">
              <thead>
                <tr>
                  <th>Nombre</th>
                  <th>Institución</th>
                  <th>Tipo / URL</th>
                  <th>Última Ejecución</th>
                  <th>Estado</th>
                  <th className="text-right">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {fuentes.length === 0 ? (
                  <tr>
                    <td colSpan="5" className="text-center py-8 text-neutral-500">No hay fuentes configuradas.</td>
                  </tr>
                ) : (
                  fuentes.map(f => (
                    <tr key={f.id}>
                      <td className="font-bold">{f.nombre}</td>
                      <td className="text-sm">{f.institucion_id || 'Global'}</td>
                      <td>
                        <span className="badge badge-neutral mr-2">{f.tipo.toUpperCase()}</span>
                        <a href={f.url} target="_blank" rel="noreferrer" className="text-primary-600 text-sm inline-flex items-center gap-1 hover:underline">
                          <Link2 size={12} /> URL
                        </a>
                      </td>
                      <td className="text-sm">
                        {f.ultima_ejecucion_en ? new Date(f.ultima_ejecucion_en).toLocaleString() : 'Nunca'}
                      </td>
                      <td>
                        <span className={`badge ${f.estado === 'activa' ? 'badge-success' : f.estado === 'error' ? 'badge-error' : 'badge-neutral'}`}>
                          {f.estado.toUpperCase()}
                        </span>
                      </td>
                      <td className="text-right">
                        <div className="flex gap-2 justify-end">
                          <button 
                            className="btn btn-sm btn-outline text-blue-600 hover:bg-blue-50 border-transparent"
                            onClick={() => openEditFuente(f)}
                          >
                            Editar
                          </button>
                          <button 
                            className="btn btn-sm text-red-600 hover:bg-red-50 border border-transparent"
                            onClick={() => handleDeleteFuente(f.id)}
                          >
                            Eliminar
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        )}

        {activeTab === 'candidatos' && (
          <div className="table-responsive">
            <table className="table">
              <thead>
                <tr>
                  <th>Confianza</th>
                  <th>Trámite Sugerido</th>
                  <th>Estado</th>
                  <th>Datos Extraídos (Preview)</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {candidatos.length === 0 ? (
                  <tr>
                    <td colSpan="5" className="text-center py-8 text-neutral-500">No hay candidatos para revisar.</td>
                  </tr>
                ) : (
                  candidatos.map(c => (
                    <tr key={c.id}>
                      <td>
                        <div className="flex items-center gap-2">
                          <div className="w-16 h-2 bg-neutral-200 rounded-full overflow-hidden">
                            <div 
                              className={`h-full ${c.confianza > 0.8 ? 'bg-green-500' : c.confianza > 0.5 ? 'bg-yellow-500' : 'bg-red-500'}`}
                              style={{ width: `${c.confianza * 100}%` }}
                            ></div>
                          </div>
                          <span className="text-xs font-bold">{Math.round(c.confianza * 100)}%</span>
                        </div>
                      </td>
                      <td className="font-mono text-sm">{c.tramite_id_sugerido || 'Desconocido'}</td>
                      <td>
                        <span className={`badge ${c.estado === 'pendiente' ? 'badge-warning' : c.estado === 'aceptado' ? 'badge-success' : 'badge-neutral'}`}>
                          {c.estado.toUpperCase()}
                        </span>
                      </td>
                      <td className="text-xs text-neutral-600 max-w-xs truncate">
                        {JSON.stringify(c.datos_extraidos).substring(0, 100)}...
                      </td>
                      <td>
                        <button className="btn btn-sm btn-outline" onClick={() => setSelectedCandidato(c)}>
                          Revisar
                        </button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {selectedCandidato && (
        <Modal isOpen={true} onClose={() => setSelectedCandidato(null)} title="Revisión de Candidato de Ingesta">
          <div className="flex-col gap-4">
            <div className="grid grid-cols-2 gap-4 text-sm bg-neutral-50 p-4 rounded-md">
              <div><span className="text-neutral-500 block">ID Candidato:</span> <span className="font-mono">{selectedCandidato.id}</span></div>
              <div><span className="text-neutral-500 block">Trámite Sugerido:</span> {selectedCandidato.tramite_id_sugerido || 'Ninguno (Creación nueva)'}</div>
              <div><span className="text-neutral-500 block">Nivel de Confianza de la IA:</span> {Math.round(selectedCandidato.confianza * 100)}%</div>
              <div><span className="text-neutral-500 block">Estado Actual:</span> {selectedCandidato.estado}</div>
            </div>

            <div>
              <h4 className="font-bold text-xs text-neutral-500 uppercase mb-2">Datos Extraídos Estructurados (JSON)</h4>
              <pre className="bg-neutral-900 text-neutral-100 p-4 rounded text-xs overflow-auto max-h-64">
                {JSON.stringify(selectedCandidato.datos_extraidos, null, 2)}
              </pre>
            </div>

            {selectedCandidato.estado === 'pendiente' ? (
              <div className="flex gap-4 mt-6">
                <button 
                  className="btn flex-1 bg-green-600 hover:bg-green-700 text-white" 
                  disabled={actionLoading}
                  onClick={() => handleUpdateCandidato('aceptado')}
                >
                  <Check size={18} /> Aceptar (Convertir en Borrador)
                </button>
                <button 
                  className="btn flex-1 bg-neutral-200 hover:bg-neutral-300 text-neutral-800" 
                  disabled={actionLoading}
                  onClick={() => handleUpdateCandidato('rechazado')}
                >
                  <X size={18} /> Rechazar
                </button>
              </div>
            ) : (
              <div className="mt-6 text-center text-sm font-medium text-neutral-500">
                Este candidato ya ha sido {selectedCandidato.estado}
              </div>
            )}
          </div>
        </Modal>
      )}

      <Modal isOpen={isFuenteModalOpen} onClose={() => setIsFuenteModalOpen(false)} title="Nueva Fuente de Ingesta">
        <form onSubmit={handleCreateFuente} className="flex-col">
          <div className="form-group">
            <label>Nombre Descriptivo</label>
            <input required className="input-field" value={fuenteForm.nombre} onChange={e => setFuenteForm({...fuenteForm, nombre: e.target.value})} placeholder="Ej. Portal SEGIP Trámites" />
          </div>
          
          <div className="form-group">
            <label>Tipo de Fuente</label>
            <select required className="input-field" value={fuenteForm.tipo} onChange={e => setFuenteForm({...fuenteForm, tipo: e.target.value})}>
              <option value="web">Página Web (Scraping)</option>
              <option value="api">API / JSON</option>
              <option value="documento">Documento (PDF/Word)</option>
            </select>
          </div>

          <div className="form-group">
            <label>URL / Endpoint</label>
            <input required type="url" className="input-field" value={fuenteForm.url} onChange={e => setFuenteForm({...fuenteForm, url: e.target.value})} placeholder="https://..." />
          </div>

          <div className="form-group">
            <label>ID Institución (Opcional)</label>
            <input type="number" className="input-field" value={fuenteForm.institucion_id} onChange={e => setFuenteForm({...fuenteForm, institucion_id: e.target.value})} placeholder="Ej. 1" />
          </div>

          <div className="form-group">
            <label>Expresión Cron (Opcional)</label>
            <input type="text" className="input-field" value={fuenteForm.frecuencia_cron} onChange={e => setFuenteForm({...fuenteForm, frecuencia_cron: e.target.value})} placeholder="0 0 * * *" />
          </div>

          {isEditingFuente && (
            <div className="form-group">
              <label>Estado</label>
              <select className="input-field" value={fuenteForm.estado} onChange={e => setFuenteForm({...fuenteForm, estado: e.target.value})}>
                <option value="activa">Activa</option>
                <option value="pausada">Pausada</option>
                <option value="error">Error</option>
              </select>
            </div>
          )}

          <div className="flex justify-between mt-6">
            <button type="button" className="btn btn-outline" onClick={() => setIsFuenteModalOpen(false)}>Cancelar</button>
            <button type="submit" className="btn btn-primary" disabled={fuenteLoading}>
              {fuenteLoading ? 'Guardando...' : (isEditingFuente ? 'Guardar Cambios' : 'Crear Fuente')}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
