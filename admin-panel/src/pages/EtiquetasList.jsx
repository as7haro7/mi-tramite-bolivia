import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { Plus, Activity } from 'lucide-react';
import { Modal } from '../components/Modal';

export default function EtiquetasList() {
  const [etiquetas, setEtiquetas] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editingItem, setEditingItem] = useState(null);

  const loadData = async () => {
    try {
      setLoading(true);
      const res = await api.get('/admin/etiquetas');
      setEtiquetas(res.datos || []);
    } catch (err) {
      setError('Error al cargar etiquetas');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = {
      slug: formData.get('slug'),
      nombre: formData.get('nombre'),
    };

    try {
      if (editingItem) {
        await api.put(`/admin/etiquetas/${editingItem.id}`, data);
      } else {
        await api.post('/admin/etiquetas', data);
      }
      setShowModal(false);
      loadData();
    } catch (err) {
      alert(err.message || 'Error al guardar');
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('¿Seguro que deseas ELIMINAR esta etiqueta permanentemente?')) return;
    try {
      await api.delete(`/admin/etiquetas/${id}`);
      loadData();
    } catch (err) {
      alert(err.message || 'Error al eliminar');
    }
  };

  const openNew = () => {
    setEditingItem(null);
    setShowModal(true);
  };

  const openEdit = (item) => {
    setEditingItem(item);
    setShowModal(true);
  };

  if (loading && etiquetas.length === 0) return <div className="page-loader"><Activity className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;

  return (
    <div className="etiquetas-list">
      <header className="page-header flex justify-between items-center">
        <div>
          <h1 className="page-title">Etiquetas</h1>
          <p className="page-subtitle">Gestión de etiquetas libres para trámites.</p>
        </div>
        <button className="btn btn-primary" onClick={openNew}>
          <Plus size={18} /> Nueva Etiqueta
        </button>
      </header>

      <div className="card mt-6">
        <div className="table-responsive">
          <table className="table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Nombre</th>
                <th>Slug</th>
                <th className="text-right">Acciones</th>
              </tr>
            </thead>
          <tbody className="divide-y divide-neutral-200">
            {etiquetas.length === 0 ? (
              <tr>
                <td colSpan="4" className="text-center py-12 text-neutral-500 text-sm">
                  No hay etiquetas registradas.
                </td>
              </tr>
            ) : (
              etiquetas.map((e) => (
                <tr key={e.id}>
                  <td className="text-sm text-neutral-500 font-mono">{e.id}</td>
                  <td className="font-bold">{e.nombre}</td>
                  <td className="text-sm text-neutral-500">{e.slug}</td>
                  <td className="text-right">
                    <div className="flex gap-2 justify-end">
                      <button className="btn btn-sm btn-outline text-blue-600 hover:bg-blue-50 border-transparent" onClick={() => openEdit(e)}>
                        Editar
                      </button>
                      <button className="btn btn-sm text-red-600 hover:bg-red-50 border border-transparent" onClick={() => handleDelete(e.id)}>
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
      </div>

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={editingItem ? 'Editar Etiqueta' : 'Nueva Etiqueta'}>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Slug</label>
            <input
              type="text"
              name="slug"
              defaultValue={editingItem?.slug || ''}
              required
              className="input-field"
              placeholder="Ej. urgente"
            />
          </div>
          <div className="form-group">
            <label>Nombre</label>
            <input
              type="text"
              name="nombre"
              defaultValue={editingItem?.nombre || ''}
              required
              className="input-field"
            />
          </div>

          <div className="flex justify-between mt-6">
            <button type="button" className="btn btn-outline" onClick={() => setShowModal(false)}>Cancelar</button>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {editingItem ? 'Guardar Cambios' : 'Crear Etiqueta'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
