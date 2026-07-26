import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../lib/api';
import { FileText, Activity, Plus } from 'lucide-react';
import { Modal } from '../components/Modal';

export function TramitesList() {
  const [tramites, setTramites] = useState([]);
  const [instituciones, setInstituciones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const [categorias, setCategorias] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [form, setForm] = useState({ id: null, codigo: '', slug: '', institucion_id: '', categoria_id: '', alcance: 'nacional', codigo_oficial: '' });
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    cargarDatos();
  }, []);

  const cargarDatos = async () => {
    try {
      setLoading(true);
      const [tResp, iResp, cResp] = await Promise.all([
        api.get('/admin/tramites'),
        api.get('/admin/instituciones'),
        api.get('/admin/categorias')
      ]);
      setTramites(tResp.datos || []);
      setInstituciones(iResp.datos || []);
      setCategorias(cResp.datos || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setFormLoading(true);
    try {
      const payload = {
        ...form,
        institucion_id: parseInt(form.institucion_id, 10),
        categoria_id: form.categoria_id ? parseInt(form.categoria_id, 10) : undefined
      };
      
      if (isEditing) {
        await api.put(`/admin/tramites/${form.id}`, payload);
        setIsModalOpen(false);
        cargarDatos();
      } else {
        const resp = await api.post('/admin/tramites', payload);
        setIsModalOpen(false);
        navigate(`/tramites/${resp.id}`);
      }
    } catch (err) {
      alert(err.data?.message || 'Error al guardar trámite');
    } finally {
      setFormLoading(false);
    }
  };

  const handleEditClick = (t) => {
    setIsEditing(true);
    setForm({
      id: t.id,
      codigo: t.codigo,
      slug: t.slug,
      institucion_id: t.institucion_id || '',
      categoria_id: t.categoria_id || '',
      alcance: t.alcance || 'nacional',
      codigo_oficial: t.codigo_oficial || ''
    });
    setIsModalOpen(true);
  };

  const handleDeleteClick = async (id) => {
    if (!confirm('¿Seguro que deseas inactivar este trámite?')) return;
    try {
      await api.delete(`/admin/tramites/${id}`);
      cargarDatos();
    } catch (err) {
      alert(err.data?.message || 'Error al inactivar');
    }
  };

  const openNewModal = () => {
    setIsEditing(false);
    setForm({ id: null, codigo: '', slug: '', institucion_id: '', categoria_id: '', alcance: 'nacional', codigo_oficial: '' });
    setIsModalOpen(true);
  };

  const badgeClass = (estadoEditorial) => {
    switch(estadoEditorial) {
      case 'borrador': return 'badge-neutral';
      case 'en_revision': return 'badge-warning';
      case 'publicada': return 'badge-success';
      case 'rechazada': return 'badge-error';
      default: return 'badge-neutral';
    }
  };

  if (loading && tramites.length === 0) return <div className="page-loader"><Activity className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;

  return (
    <div className="tramites-list">
      <header className="page-header flex justify-between items-center">
        <div>
          <h1 className="page-title">Catálogo de Trámites</h1>
          <p className="page-subtitle">Gestión de trámites y flujo editorial</p>
        </div>
        <button className="btn btn-primary" onClick={openNewModal}>
          <Plus size={18} /> Nuevo Trámite
        </button>
      </header>

      <div className="card mt-6">
        <div className="table-responsive">
          <table className="table">
            <thead>
              <tr>
                <th>Código</th>
                <th>Institución</th>
                <th>Versión Actual</th>
                <th>Estado Editorial</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {tramites.length === 0 ? (
                <tr>
                  <td colSpan="5" className="text-center py-8 text-neutral-500">
                    No hay trámites registrados.
                  </td>
                </tr>
              ) : (
                tramites.map(t => (
                  <tr key={t.id}>
                    <td className="font-mono text-sm">{t.codigo}</td>
                    <td>{t.institucion}</td>
                    <td>
                      {t.numero_version ? (
                        <div>
                          <span className="font-medium">v{t.numero_version}</span>
                          <span className="text-xs text-neutral-500 block truncate max-w-xs" title={t.titulo}>
                            {t.titulo}
                          </span>
                        </div>
                      ) : (
                        <span className="text-neutral-400 italic">Sin versiones</span>
                      )}
                    </td>
                    <td>
                      {t.estado_editorial ? (
                        <span className={`badge ${badgeClass(t.estado_editorial)}`}>
                          {t.estado_editorial.replace('_', ' ').toUpperCase()}
                        </span>
                      ) : (
                        <span className="badge badge-neutral">SIN VERSIÓN</span>
                      )}
                    </td>
                    <td className="text-right">
                      <div className="flex gap-2 justify-end">
                        <button 
                          className="btn btn-sm btn-outline"
                          onClick={() => navigate(`/tramites/${t.id}`)}
                        >
                          <FileText size={16} /> Ver Detalles
                        </button>
                        <button 
                          className="btn btn-sm btn-outline text-blue-600 hover:bg-blue-50"
                          onClick={() => handleEditClick(t)}
                          title="Editar Base"
                        >
                          Editar
                        </button>
                        {t.estado !== 'retirado' && (
                          <button 
                            className="btn btn-sm text-red-600 hover:bg-red-50 border border-transparent"
                            onClick={() => handleDeleteClick(t.id)}
                            title="Retirar Trámite"
                          >
                            Eliminar
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="Nuevo Trámite (Base)">
        <form onSubmit={handleSubmit}>
          <div className="alert bg-blue-50 text-blue-800 border-blue-200 mb-4">
            Esto crea el registro base del trámite. La información detallada (título, requisitos) se ingresa al crear una "Versión".
          </div>
          <div className="form-group">
            <label>Código Oficial (ej. SEGIP-01)</label>
            <input required className="input-field" value={form.codigo} onChange={e => setForm({...form, codigo: e.target.value})} />
          </div>
          <div className="form-group">
            <label>URL Slug (ej. emision-carnet-identidad)</label>
            <input required className="input-field" value={form.slug} onChange={e => setForm({...form, slug: e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, '-')})} />
          </div>
          <div className="form-group">
            <label>Institución Responsable</label>
            <select required className="input-field" value={form.institucion_id} onChange={e => setForm({...form, institucion_id: e.target.value})}>
              <option value="">Seleccione una institución...</option>
              {instituciones.filter(i => i.estado === 'activo').map(i => (
                <option key={i.id} value={i.id}>{i.sigla} - {i.nombre}</option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label>Categoría</label>
            <select className="input-field" value={form.categoria_id} onChange={e => setForm({...form, categoria_id: e.target.value})}>
              <option value="">Sin Categoría</option>
              {categorias.filter(c => c.activa).map(c => (
                <option key={c.id} value={c.id}>{c.nombre}</option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label>Alcance</label>
            <select required className="input-field" value={form.alcance} onChange={e => setForm({...form, alcance: e.target.value})}>
              <option value="nacional">Nacional</option>
              <option value="departamental">Departamental</option>
              <option value="municipal">Municipal</option>
            </select>
          </div>
          
          <div className="flex justify-between mt-6">
            <button type="button" className="btn btn-outline" onClick={() => setIsModalOpen(false)}>Cancelar</button>
            <button type="submit" className="btn btn-primary" disabled={formLoading}>
              {formLoading ? 'Guardando...' : (isEditing ? 'Guardar Cambios' : 'Crear y Continuar')}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
