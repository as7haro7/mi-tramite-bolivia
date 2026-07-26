import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { Building2, Activity, Globe, Plus } from 'lucide-react';
import { Modal } from '../components/Modal';

export function InstitucionesList() {
  const [instituciones, setInstituciones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form, setForm] = useState({ codigo: '', sigla: '', nombre: '', tipo: 'publica', sitio_web: '' });
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    cargarInstituciones();
  }, []);

  const cargarInstituciones = async () => {
    try {
      setLoading(true);
      const resp = await api.get('/admin/instituciones');
      setInstituciones(resp.datos || []);
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
      if (form.id) {
        await api.put(`/admin/instituciones/${form.id}`, form);
      } else {
        await api.post('/admin/instituciones', form);
      }
      setIsModalOpen(false);
      setForm({ id: null, codigo: '', sigla: '', nombre: '', tipo: 'publica', sitio_web: '' });
      cargarInstituciones();
    } catch (err) {
      alert(err.data?.message || 'Error al guardar institución');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('¿Seguro que deseas inactivar esta institución?')) return;
    try {
      await api.delete(`/admin/instituciones/${id}`);
      cargarInstituciones();
    } catch (err) {
      alert(err.data?.message || 'Error al inactivar');
    }
  };

  const openEdit = (inst) => {
    setForm({
      id: inst.id,
      codigo: inst.codigo,
      sigla: inst.sigla,
      nombre: inst.nombre,
      tipo: inst.tipo || 'publica',
      sitio_web: inst.sitio_web || ''
    });
    setIsModalOpen(true);
  };

  const openCreate = () => {
    setForm({ id: null, codigo: '', sigla: '', nombre: '', tipo: 'publica', sitio_web: '' });
    setIsModalOpen(true);
  };

  if (loading && instituciones.length === 0) return <div className="page-loader"><Activity className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;

  return (
    <div className="instituciones-list">
      <header className="page-header flex justify-between items-center">
        <div>
          <h1 className="page-title">Instituciones</h1>
          <p className="page-subtitle">Gestión de catálogo de instituciones</p>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          <Plus size={18} /> Nueva Institución
        </button>
      </header>

      <div className="card mt-6">
        <div className="table-responsive">
          <table className="table">
            <thead>
              <tr>
                <th>Código / Sigla</th>
                <th>Nombre</th>
                <th>Tipo</th>
                <th>Sitio Web</th>
                <th>Estado</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {instituciones.length === 0 ? (
                <tr>
                  <td colSpan="6" className="text-center py-8 text-neutral-500">
                    No hay instituciones registradas.
                  </td>
                </tr>
              ) : (
                instituciones.map(i => (
                  <tr key={i.id}>
                    <td>
                      <div className="font-mono text-sm">{i.codigo}</div>
                      <div className="font-bold text-xs text-primary">{i.sigla}</div>
                    </td>
                    <td className="font-medium">{i.nombre}</td>
                    <td className="capitalize">{i.tipo}</td>
                    <td>
                      {i.sitio_web ? (
                        <a href={i.sitio_web} target="_blank" rel="noreferrer" className="text-primary hover:underline inline-flex items-center gap-1">
                          <Globe size={14} /> Link
                        </a>
                      ) : (
                        <span className="text-neutral-400">-</span>
                      )}
                    </td>
                    <td>
                      <span className={`badge ${i.estado === 'activo' ? 'badge-success' : 'badge-neutral'}`}>
                        {i.estado.toUpperCase()}
                      </span>
                    </td>
                    <td>
                      <div className="flex gap-2 justify-end">
                        <button className="btn btn-sm btn-outline text-blue-600 hover:bg-blue-50 border-transparent" onClick={() => openEdit(i)}>
                          Editar
                        </button>
                        {i.estado !== 'inactiva' && (
                          <button 
                            className="btn btn-sm text-red-600 hover:bg-red-50 border border-transparent"
                            onClick={() => handleDelete(i.id)}
                            title="Inactivar Institución"
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

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={form.id ? 'Editar Institución' : 'Nueva Institución'}>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Código (URL safe)</label>
            <input required className="input-field" value={form.codigo} onChange={e => setForm({...form, codigo: e.target.value})} placeholder="ej. sin" />
          </div>
          <div className="form-group">
            <label>Sigla</label>
            <input required className="input-field" value={form.sigla} onChange={e => setForm({...form, sigla: e.target.value})} placeholder="ej. SIN" />
          </div>
          <div className="form-group">
            <label>Nombre Completo</label>
            <input required className="input-field" value={form.nombre} onChange={e => setForm({...form, nombre: e.target.value})} placeholder="ej. Servicio de Impuestos Nacionales" />
          </div>
          <div className="form-group">
            <label>Sitio Web (opcional)</label>
            <input type="url" className="input-field" value={form.sitio_web} onChange={e => setForm({...form, sitio_web: e.target.value})} placeholder="https://" />
          </div>
          
          <div className="flex justify-between mt-6">
            <button type="button" className="btn btn-outline" onClick={() => setIsModalOpen(false)}>Cancelar</button>
            <button type="submit" className="btn btn-primary" disabled={formLoading}>
              {formLoading ? 'Guardando...' : 'Guardar Institución'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
