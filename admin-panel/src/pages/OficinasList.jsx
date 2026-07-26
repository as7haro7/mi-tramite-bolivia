import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { Activity, Plus, MapPin } from 'lucide-react';
import { Modal } from '../components/Modal';

export function OficinasList() {
  const [oficinas, setOficinas] = useState([]);
  const [instituciones, setInstituciones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form, setForm] = useState({ 
    institucion_id: '', 
    codigo: '', 
    nombre: '', 
    tipo: 'oficina', 
    direccion: '', 
    zona_horaria: 'America/La_Paz' 
  });
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    cargarDatos();
  }, []);

  const cargarDatos = async () => {
    try {
      setLoading(true);
      const [oResp, iResp] = await Promise.all([
        api.get('/admin/oficinas'),
        api.get('/admin/instituciones')
      ]);
      setOficinas(oResp.datos || []);
      setInstituciones(iResp.datos || []);
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
        institucion_id: parseInt(form.institucion_id, 10)
      };
      
      if (form.id) {
        await api.put(`/admin/oficinas/${form.id}`, payload);
      } else {
        await api.post('/admin/oficinas', payload);
      }
      setIsModalOpen(false);
      setForm({ id: null, institucion_id: '', codigo: '', nombre: '', tipo: 'oficina', direccion: '', zona_horaria: 'America/La_Paz' });
      cargarDatos();
    } catch (err) {
      alert(err.data?.message || 'Error al guardar oficina');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('¿Seguro que deseas inactivar esta oficina?')) return;
    try {
      await api.delete(`/admin/oficinas/${id}`);
      cargarDatos();
    } catch (err) {
      alert(err.data?.message || 'Error al inactivar');
    }
  };

  const openCreate = () => {
    setForm({ id: null, institucion_id: '', codigo: '', nombre: '', tipo: 'oficina', direccion: '', zona_horaria: 'America/La_Paz' });
    setIsModalOpen(true);
  };

  const openEdit = (oficina) => {
    setForm({
      id: oficina.id,
      institucion_id: oficina.institucion_id || '',
      codigo: oficina.codigo || '',
      nombre: oficina.nombre || '',
      tipo: oficina.tipo || 'oficina',
      direccion: oficina.direccion || '',
      zona_horaria: oficina.zona_horaria || 'America/La_Paz'
    });
    setIsModalOpen(true);
  };

  if (loading && oficinas.length === 0) return <div className="page-loader"><Activity className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;

  return (
    <div className="oficinas-list">
      <header className="page-header flex justify-between items-center">
        <div>
          <h1 className="page-title">Puntos de Atención</h1>
          <p className="page-subtitle">Gestión de oficinas físicas y virtuales</p>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          <Plus size={18} /> Nueva Oficina
        </button>
      </header>

      <div className="card mt-6">
        <div className="table-responsive">
          <table className="table">
            <thead>
              <tr>
                <th>Código</th>
                <th>Institución</th>
                <th>Nombre / Tipo</th>
                <th>Ubicación</th>
                <th>Estado</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {oficinas.length === 0 ? (
                <tr>
                  <td colSpan="6" className="text-center py-8 text-neutral-500">
                    No hay oficinas registradas.
                  </td>
                </tr>
              ) : (
                oficinas.map(o => (
                  <tr key={o.id}>
                    <td className="font-mono text-sm">{o.codigo}</td>
                    <td className="font-medium text-sm">{o.institucion}</td>
                    <td>
                      <div className="font-bold">{o.nombre}</div>
                      <div className="text-xs text-neutral-500 capitalize">{o.tipo.replace('_', ' ')}</div>
                    </td>
                    <td>
                      {o.direccion ? (
                        <div className="inline-flex items-center gap-1 text-sm text-neutral-600">
                          <MapPin size={14} /> {o.direccion}
                        </div>
                      ) : <span className="text-neutral-400">-</span>}
                    </td>
                    <td>
                      <span className={`badge ${o.estado === 'activa' ? 'badge-success' : 'badge-neutral'}`}>
                        {o.estado.toUpperCase()}
                      </span>
                    </td>
                    <td>
                      <div className="flex gap-2 justify-end">
                        <button className="btn btn-sm btn-outline text-blue-600 hover:bg-blue-50 border-transparent" onClick={() => openEdit(o)}>
                          Editar
                        </button>
                        {o.estado !== 'inactiva' && (
                          <button 
                            className="btn btn-sm text-red-600 hover:bg-red-50 border border-transparent"
                            onClick={() => handleDelete(o.id)}
                            title="Inactivar Oficina"
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

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={form.id ? 'Editar Oficina' : 'Nueva Oficina / Punto'}>
        <form onSubmit={handleSubmit} className="flex-col">
          <div className="form-group">
            <label>Institución a la que pertenece</label>
            <select required className="input-field" value={form.institucion_id} onChange={e => setForm({...form, institucion_id: e.target.value})}>
              <option value="">Seleccione institución...</option>
              {instituciones.filter(i => i.estado === 'activo').map(i => (
                <option key={i.id} value={i.id}>{i.sigla} - {i.nombre}</option>
              ))}
            </select>
          </div>
          
          <div className="flex-row">
            <div className="form-group flex-1">
              <label>Código Interno</label>
              <input required className="input-field" value={form.codigo} onChange={e => setForm({...form, codigo: e.target.value})} placeholder="Ej. OFI-CEN-01" />
            </div>
            <div className="form-group flex-1">
              <label>Tipo de Atención</label>
              <select required className="input-field" value={form.tipo} onChange={e => setForm({...form, tipo: e.target.value})}>
                <option value="oficina">Oficina Central/Regional</option>
                <option value="ventanilla">Ventanilla Única</option>
                <option value="brigada_movil">Brigada Móvil</option>
                <option value="punto_pago">Punto de Pago</option>
                <option value="virtual">Oficina Virtual / Portal</option>
              </select>
            </div>
          </div>

          <div className="form-group">
            <label>Nombre Descriptivo</label>
            <input required className="input-field" value={form.nombre} onChange={e => setForm({...form, nombre: e.target.value})} placeholder="Ej. Oficina Central La Paz" />
          </div>
          
          <div className="form-group">
            <label>Dirección Física (opcional)</label>
            <textarea className="input-field" style={{minHeight: '60px'}} value={form.direccion} onChange={e => setForm({...form, direccion: e.target.value})} placeholder="Av. Principal #123..." />
          </div>

          <div className="flex justify-between mt-6">
            <button type="button" className="btn btn-outline" onClick={() => setIsModalOpen(false)}>Cancelar</button>
            <button type="submit" className="btn btn-primary" disabled={formLoading}>
              {formLoading ? 'Guardando...' : 'Guardar Oficina'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
