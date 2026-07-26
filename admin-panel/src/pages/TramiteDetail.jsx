import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { api } from '../lib/api';
import { getUsuario } from '../lib/auth';
import { Activity, ArrowLeft, Save, Send, CheckCircle, XCircle, Globe, AlertCircle } from 'lucide-react';
import { Modal } from '../components/Modal';

export function TramiteDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const usuario = getUsuario();
  
  const [tramite, setTramite] = useState(null);
  const [versiones, setVersiones] = useState([]);
  const [versionActiva, setVersionActiva] = useState(null);
  
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  
  // Editor form
  const [form, setForm] = useState(null);
  const [saving, setSaving] = useState(false);

  // Rejection modal
  const [rechazoModal, setRechazoModal] = useState(false);
  const [observaciones, setObservaciones] = useState('');

  useEffect(() => {
    cargarDetalles();
  }, [id]);

  const cargarDetalles = async () => {
    try {
      setLoading(true);
      // Cargar trámite y sus detalles
      // En este caso llamamos al API público o admin según disponibilidad
      // Usar el nuevo endpoint admin que trae todo el detalle
      const t = await api.get(`/admin/tramites/${id}`);
      
      if (!t) throw new Error('Trámite no encontrado');
      
      setTramite(t);
      
      // Si tiene versión actual, creamos su form
      if (t.version_id) {
        setVersionActiva({
          id: t.version_id,
          numero: t.numero_version,
          estado_editorial: t.estado_editorial,
          titulo: t.titulo || '',
          resumen: t.resumen || '',
          descripcion: t.descripcion || '',
          publico_objetivo: t.publico_objetivo || '',
          resultado_esperado: t.resultado_esperado || '',
          advertencias: t.advertencias || '',
          plazo_texto: t.plazo_texto || '',
          url_inicio: t.url_inicio || ''
        });
        
        setForm({
          titulo: t.titulo || '',
          resumen: t.resumen || '',
          descripcion: t.descripcion || '',
          publico_objetivo: t.publico_objetivo || '',
          resultado_esperado: t.resultado_esperado || '',
          advertencias: t.advertencias || '',
          plazo_texto: t.plazo_texto || '',
          url_inicio: t.url_inicio || ''
        });
      }

    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCrearVersion = async () => {
    try {
      setSaving(true);
      const resp = await api.post(`/admin/tramites/${id}/versiones`, {
        titulo: 'Nuevo Trámite',
        resumen: 'Resumen inicial',
        descripcion: 'Descripción inicial'
      });
      cargarDetalles();
    } catch (err) {
      alert(err.data?.message || 'Error al crear versión');
    } finally {
      setSaving(false);
    }
  };

  const handleGuardarCambios = async (e) => {
    e.preventDefault();
    if (!versionActiva) return;
    
    try {
      setSaving(true);
      await api.put(`/admin/versiones/${versionActiva.id}`, form);
      alert('Cambios guardados exitosamente');
      cargarDetalles();
    } catch (err) {
      alert(err.data?.message || 'Error al guardar');
    } finally {
      setSaving(false);
    }
  };

  const cambiarEstado = async (accion, payload = {}) => {
    if (!versionActiva) return;
    try {
      setSaving(true);
      await api.post(`/admin/versiones/${versionActiva.id}/${accion}`, payload);
      setRechazoModal(false);
      cargarDetalles();
    } catch (err) {
      alert(err.data?.message || 'Error al cambiar estado');
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <div className="page-loader"><Activity className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;
  if (!tramite) return null;

  const isBorrador = versionActiva?.estado_editorial === 'borrador';
  const isEnRevision = versionActiva?.estado_editorial === 'en_revision';
  const isPublicada = versionActiva?.estado_editorial === 'publicada';
  
  // Un superadmin o editor jefe podría tener permiso para aprobar/publicar
  const canApprove = usuario?.roles?.includes('superadmin') || usuario?.roles?.includes('editor_jefe');

  return (
    <div className="tramite-detail">
      <header className="page-header flex justify-between items-center mb-6">
        <div className="flex items-center gap-4">
          <button className="btn btn-outline" onClick={() => navigate('/tramites')}>
            <ArrowLeft size={18} />
          </button>
          <div>
            <h1 className="page-title">{tramite.codigo} - {tramite.institucion}</h1>
            <p className="page-subtitle">Gestión Editorial de Trámite</p>
          </div>
        </div>
      </header>

      {!versionActiva ? (
        <div className="card p-12 text-center flex-col items-center">
          <AlertCircle size={48} className="text-neutral-400 mx-auto mb-4" />
          <h2 className="text-xl font-bold mb-2">Sin versiones activas</h2>
          <p className="text-neutral-500 mb-6">Este trámite acaba de ser creado y no tiene contenido.</p>
          <button className="btn btn-primary" onClick={handleCrearVersion} disabled={saving}>
            {saving ? 'Creando...' : 'Crear Versión 1 (Borrador)'}
          </button>
        </div>
      ) : (
        <div className="flex-row items-start">
          {/* Panel Izquierdo: Formulario */}
          <div className="card flex-1 p-6">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-lg font-bold">Contenido de la Versión {versionActiva.numero}</h2>
              <span className={`badge ${isBorrador ? 'badge-neutral' : isEnRevision ? 'badge-warning' : 'badge-success'}`}>
                {versionActiva.estado_editorial.toUpperCase()}
              </span>
            </div>

            <form onSubmit={handleGuardarCambios} className="flex-col">
              <div className="form-group">
                <label>Título del Trámite</label>
                <input className="input-field" value={form.titulo} onChange={e => setForm({...form, titulo: e.target.value})} disabled={!isBorrador} required />
              </div>

              <div className="form-group">
                <label>Resumen Breve</label>
                <textarea className="input-field" value={form.resumen} onChange={e => setForm({...form, resumen: e.target.value})} disabled={!isBorrador} required />
              </div>

              <div className="form-group">
                <label>Descripción Completa</label>
                <textarea className="input-field" style={{minHeight: '150px'}} value={form.descripcion} onChange={e => setForm({...form, descripcion: e.target.value})} disabled={!isBorrador} required />
              </div>

              <div className="flex-row">
                <div className="form-group flex-1">
                  <label>Público Objetivo</label>
                  <input className="input-field" value={form.publico_objetivo} onChange={e => setForm({...form, publico_objetivo: e.target.value})} disabled={!isBorrador} />
                </div>
                <div className="form-group flex-1">
                  <label>Resultado Esperado</label>
                  <input className="input-field" value={form.resultado_esperado} onChange={e => setForm({...form, resultado_esperado: e.target.value})} disabled={!isBorrador} />
                </div>
              </div>

              <div className="form-group">
                <label>Advertencias Importantes</label>
                <textarea className="input-field" value={form.advertencias} onChange={e => setForm({...form, advertencias: e.target.value})} disabled={!isBorrador} />
              </div>

              <div className="flex-row">
                <div className="form-group flex-1">
                  <label>Plazo de Entrega (Texto)</label>
                  <input className="input-field" value={form.plazo_texto} onChange={e => setForm({...form, plazo_texto: e.target.value})} disabled={!isBorrador} />
                </div>
                <div className="form-group flex-1">
                  <label>URL Trámite en Línea (opcional)</label>
                  <input type="url" className="input-field" value={form.url_inicio} onChange={e => setForm({...form, url_inicio: e.target.value})} disabled={!isBorrador} />
                </div>
              </div>

              {isBorrador && (
                <div className="mt-4 border-t pt-4">
                  <button type="submit" className="btn btn-primary" disabled={saving}>
                    <Save size={18} /> {saving ? 'Guardando...' : 'Guardar Cambios'}
                  </button>
                </div>
              )}
            </form>
          </div>

          {/* Panel Derecho: Controles Editoriales */}
          <div className="card p-6" style={{width: '350px'}}>
            <h3 className="font-bold mb-4">Acciones Editoriales</h3>
            
            <div className="flex-col">
              {isBorrador && (
                <>
                  <div className="alert bg-blue-50 text-blue-800 border-blue-200">
                    Estás editando un borrador. Cuando termines, envíalo a revisión.
                  </div>
                  <button className="btn btn-outline w-full" onClick={() => cambiarEstado('enviar-revision')} disabled={saving}>
                    <Send size={18} /> Solicitar Revisión
                  </button>
                </>
              )}

              {isEnRevision && (
                <>
                  <div className="alert bg-yellow-50 text-yellow-800 border-yellow-200">
                    Versión bloqueada. En espera de revisión por un editor jefe.
                  </div>
                  {canApprove ? (
                    <>
                      <button className="btn btn-success w-full text-white bg-green-600 hover:bg-green-700" onClick={() => cambiarEstado('aprobar')} disabled={saving}>
                        <CheckCircle size={18} /> Aprobar Versión
                      </button>
                      <button className="btn btn-error w-full text-white bg-red-600 hover:bg-red-700" onClick={() => setRechazoModal(true)} disabled={saving}>
                        <XCircle size={18} /> Rechazar
                      </button>
                    </>
                  ) : (
                    <p className="text-sm text-neutral-500 text-center">No tienes permisos para aprobar.</p>
                  )}
                </>
              )}

              {isPublicada && (
                <>
                  <div className="alert bg-green-50 text-green-800 border-green-200">
                    Esta es la versión pública activa de la institución.
                  </div>
                  <button className="btn btn-primary w-full" onClick={handleCrearVersion} disabled={saving}>
                    Crear Nueva Versión
                  </button>
                </>
              )}

              {/* Si está aprobada pero el backend v2 asume en_revision -> publicar directamente */}
              {isEnRevision && canApprove && (
                <div className="mt-4 pt-4 border-t border-dashed">
                  <button className="btn btn-primary w-full" onClick={() => cambiarEstado('publicar')} disabled={saving}>
                    <Globe size={18} /> Publicar Directamente
                  </button>
                  <p className="text-xs text-neutral-500 mt-2 text-center">
                    Esto cerrará la versión anterior y publicará esta.
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      <Modal isOpen={rechazoModal} onClose={() => setRechazoModal(false)} title="Rechazar Versión">
        <div className="form-group">
          <label>Motivo del rechazo (Observaciones)</label>
          <textarea 
            className="input-field" 
            value={observaciones} 
            onChange={e => setObservaciones(e.target.value)} 
            placeholder="Faltan requisitos, mala redacción..." 
          />
        </div>
        <div className="flex justify-between mt-6">
          <button className="btn btn-outline" onClick={() => setRechazoModal(false)}>Cancelar</button>
          <button className="btn bg-red-600 text-white hover:bg-red-700" onClick={() => cambiarEstado('rechazar', { observaciones })}>
            Rechazar y Devolver a Borrador
          </button>
        </div>
      </Modal>

    </div>
  );
}
