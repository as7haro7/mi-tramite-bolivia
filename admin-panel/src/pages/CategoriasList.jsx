import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { Plus, FolderTree, Activity } from 'lucide-react';
import { Modal } from '../components/Modal';

export default function CategoriasList() {
  const [categorias, setCategorias] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editingItem, setEditingItem] = useState(null);

  const loadData = async () => {
    try {
      setLoading(true);
      const res = await api.get('/admin/categorias');
      setCategorias(res.datos || []);
    } catch (err) {
      setError('Error al cargar categorías');
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
      codigo: formData.get('codigo'),
      nombre: formData.get('nombre'),
      icono: formData.get('icono') || null,
      orden: parseInt(formData.get('orden') || '0', 10),
      padre_id: formData.get('padre_id') ? parseInt(formData.get('padre_id'), 10) : null,
    };

    try {
      if (editingItem) {
        // En update permitimos cambiar 'activa' (borrado lógico lo hace inactiva)
        data.activa = formData.get('activa') === 'on';
        await api.put(`/admin/categorias/${editingItem.id}`, data);
      } else {
        await api.post('/admin/categorias', data);
      }
      setShowModal(false);
      loadData();
    } catch (err) {
      alert(err.message || 'Error al guardar');
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('¿Seguro que deseas inactivar esta categoría?')) return;
    try {
      await api.delete(`/admin/categorias/${id}`);
      loadData();
    } catch (err) {
      alert(err.message || 'Error al inactivar');
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

  if (loading && categorias.length === 0) return <div className="page-loader"><Activity className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;

  return (
    <div className="categorias-list">
      <header className="page-header flex justify-between items-center">
        <div>
          <h1 className="page-title">Categorías</h1>
          <p className="page-subtitle">Gestión de categorías para clasificación de trámites.</p>
        </div>
        <button className="btn btn-primary" onClick={openNew}>
          <Plus size={18} /> Nueva Categoría
        </button>
      </header>

      <div className="card mt-6">
        <div className="table-responsive">
          <table className="table">
            <thead>
              <tr>
                <th>Orden</th>
                <th>Código / Nombre</th>
                <th>Ícono</th>
                <th>Estado</th>
                <th className="text-right">Acciones</th>
              </tr>
            </thead>
          <tbody className="divide-y divide-neutral-200">
            {categorias.length === 0 ? (
              <tr>
                <td colSpan="5" className="px-6 py-12 text-center text-neutral-500 text-sm">
                  No hay categorías registradas.
                </td>
              </tr>
            ) : (
              categorias.map((c) => (
                <tr key={c.id}>
                  <td className="text-sm font-medium">{c.orden}</td>
                  <td>
                    <div className="font-bold">{c.nombre}</div>
                    <div className="text-xs text-neutral-500">{c.codigo}</div>
                  </td>
                  <td>
                    {c.icono ? <span className="material-icons text-neutral-400 text-sm">{c.icono}</span> : '-'}
                  </td>
                  <td>
                    <span className={`badge ${c.activa ? 'badge-success' : 'badge-error'}`}>
                      {c.activa ? 'ACTIVA' : 'INACTIVA'}
                    </span>
                  </td>
                  <td className="text-right">
                    <div className="flex gap-2 justify-end">
                      <button className="btn btn-sm btn-outline text-blue-600 hover:bg-blue-50 border-transparent" onClick={() => openEdit(c)}>
                        Editar
                      </button>
                      {c.activa && (
                        <button className="btn btn-sm text-red-600 hover:bg-red-50 border border-transparent" onClick={() => handleDelete(c.id)}>
                          Inactivar
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

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={editingItem ? 'Editar Categoría' : 'Nueva Categoría'}>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Código</label>
            <input
              type="text"
              name="codigo"
              defaultValue={editingItem?.codigo || ''}
              required
              className="input-field"
              placeholder="Ej. educacion"
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
          <div className="form-group">
            <label>Ícono (Material Icon name)</label>
            <input
              type="text"
              name="icono"
              defaultValue={editingItem?.icono || ''}
              className="input-field"
              placeholder="Ej. school"
            />
          </div>
          <div className="form-group">
            <label>Orden</label>
            <input
              type="number"
              name="orden"
              defaultValue={editingItem?.orden || 0}
              className="input-field"
            />
          </div>
          <div className="form-group">
            <label>Categoría Padre (Opcional)</label>
            <select
              name="padre_id"
              defaultValue={editingItem?.padre_id || ''}
              className="input-field"
            >
              <option value="">Ninguno</option>
              {categorias.filter(c => c.id !== editingItem?.id).map(c => (
                <option key={c.id} value={c.id}>{c.nombre}</option>
              ))}
            </select>
          </div>
          
          {editingItem && (
            <div className="flex items-center gap-2 mt-4">
              <input
                id="activa"
                name="activa"
                type="checkbox"
                defaultChecked={editingItem.activa}
                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded"
              />
              <label htmlFor="activa" className="block text-sm text-neutral-900 font-medium">
                Categoría Activa
              </label>
            </div>
          )}

          <div className="flex justify-between mt-6">
            <button type="button" className="btn btn-outline" onClick={() => setShowModal(false)}>Cancelar</button>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {editingItem ? 'Guardar Cambios' : 'Crear Categoría'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
